/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 00:00:00
 * @FilePath: \go-swagger\documents\definitions.go
 * @Description: 独立子文档定义合并 - 递归合并 definitions、引用收集
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package documents

import (
	"strings"

	"github.com/kamalyes/go-swagger/constants"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// mergeDocumentDefinitions 递归合并独立文档所需 definitions
func (b *Builder) mergeDocumentDefinitions(targetDefinitions map[string]interface{}, serviceSpec map[string]interface{}, selectedPaths map[string]map[string]interface{}, serviceName string) {
	sourceDefinitions, ok := serviceSpec[constants.FieldDefs].(map[string]interface{})
	if !ok || len(sourceDefinitions) == 0 {
		return
	}

	queued := make(map[string]bool)
	processed := make(map[string]bool)
	queue := make([]string, 0)

	for _, pathItem := range selectedPaths {
		enqueueDefinitionRefs(pathItem, queued, &queue)
	}

	for len(queue) > 0 {
		definitionName := queue[0]
		queue = queue[1:]

		if processed[definitionName] {
			continue
		}
		processed[definitionName] = true

		definition, exists := sourceDefinitions[definitionName]
		if !exists {
			b.logger.Warn("独立文档引用的 definition 不存在: service=%s, definition=%s", serviceName, definitionName)
			continue
		}

		clonedDefinition := mathx.ConvertMapKeysToString(definition)
		if existingDefinition, exists := targetDefinitions[definitionName]; exists {
			b.checkDefinitionConsistency(definitionName, existingDefinition, clonedDefinition, serviceName)
		} else {
			targetDefinitions[definitionName] = clonedDefinition
		}

		enqueueDefinitionRefs(clonedDefinition, queued, &queue)
	}
}

// enqueueDefinitionRefs 从对象中提取 definition 引用并入队
func enqueueDefinitionRefs(obj interface{}, queued map[string]bool, queue *[]string) {
	refs := make(map[string]struct{})
	collectDefinitionRefs(obj, refs)

	for refName := range refs {
		if queued[refName] {
			continue
		}
		queued[refName] = true
		*queue = append(*queue, refName)
	}
}

// collectDefinitionRefs 递归收集 #/definitions 引用
func collectDefinitionRefs(obj interface{}, refs map[string]struct{}) {
	switch value := obj.(type) {
	case map[string]interface{}:
		if refName := extractDefinitionName(value); refName != "" {
			refs[refName] = struct{}{}
		}
		for _, nested := range value {
			collectDefinitionRefs(nested, refs)
		}
	case []interface{}:
		for _, nested := range value {
			collectDefinitionRefs(nested, refs)
		}
	}
}

// extractDefinitionName 提取 definition 名称
func extractDefinitionName(value map[string]interface{}) string {
	refValue, ok := value[constants.FieldRef].(string)
	if !ok || !strings.HasPrefix(refValue, constants.PathDefinitions) {
		return ""
	}

	return strings.TrimPrefix(refValue, constants.PathDefinitions)
}
