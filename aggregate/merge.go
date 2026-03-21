/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 00:00:00
 * @FilePath: \go-swagger\aggregate\merge.go
 * @Description: 聚合合并模式 - 将多个服务规范合并为统一规范
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package aggregate

import (
	"encoding/json"

	"github.com/kamalyes/go-swagger/constants"
	serrors "github.com/kamalyes/go-swagger/errors"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// mergeAllSpecs 合并所有服务规范为统一规范
func (a *Aggregator) mergeAllSpecs() error {
	a.aggregatedSpec = a.newBaseSpec()
	a.aggregatedSpec[constants.FieldTags] = make([]interface{}, 0)

	serviceNames := a.getSortedServiceNames()

	allPaths := a.aggregatedSpec[constants.FieldPaths].(map[string]interface{})
	allDefinitions := a.aggregatedSpec[constants.FieldDefs].(map[string]interface{})
	allTags := a.aggregatedSpec[constants.FieldTags].([]interface{})
	tagNames := make(map[string]bool)

	for _, serviceName := range serviceNames {
		if err := a.mergeServiceSpec(serviceName, allPaths, allDefinitions, &allTags, tagNames); err != nil {
			return err
		}
	}

	a.aggregatedSpec[constants.FieldTags] = allTags

	if err := a.fixReferences(); err != nil {
		a.logger.Error("修复引用路径失败: %v", err)
	}

	a.logger.Info("规范合并完成，路径数: %d，定义数: %d", len(allPaths), len(allDefinitions))
	return nil
}

// mergeServiceSpec 合并单个服务的规范到聚合结果
func (a *Aggregator) mergeServiceSpec(serviceName string, allPaths, allDefinitions map[string]interface{}, allTags *[]interface{}, tagNames map[string]bool) error {
	spec := a.serviceSpecs[serviceName]
	a.logger.Info("正在合并服务 %s 的规范", serviceName)

	convertedSpec := mathx.ConvertMapKeysToString(spec)
	specMap, ok := convertedSpec.(map[string]interface{})
	if !ok {
		return newMergeError(serviceName, "无法转换为map[string]interface{}")
	}

	a.mergePaths(specMap, serviceName, allPaths)
	a.mergeDefinitions(specMap, serviceName, allDefinitions)
	a.mergeServiceSpecTags(specMap, serviceName, allTags, tagNames)

	return nil
}

// mergePaths 合并服务规范中的路径
func (a *Aggregator) mergePaths(specMap map[string]interface{}, serviceName string, allPaths map[string]interface{}) {
	paths, ok := specMap[constants.FieldPaths].(map[string]interface{})
	if !ok {
		return
	}

	for path, operations := range paths {
		convertedOps := mathx.ConvertMapKeysToString(operations)
		newOps, ok := convertedOps.(map[string]interface{})
		if !ok {
			a.logger.Error("路径 %s 的操作格式不正确", path)
			continue
		}
		a.mergePathOperations(path, newOps, serviceName, allPaths)
	}
}

// mergePathOperations 合并单个路径的操作
func (a *Aggregator) mergePathOperations(path string, newOps map[string]interface{}, serviceName string, allPaths map[string]interface{}) {
	existingPath, exists := allPaths[path]
	if !exists {
		allPaths[path] = newOps
		a.logger.Debug("添加新路径: %s (来自: %s)", path, serviceName)
		return
	}

	existingOps, ok := existingPath.(map[string]interface{})
	if !ok {
		a.logger.Error("现有路径 %s 的操作格式不正确", path)
		return
	}

	mergedAny := false
	for method, op := range newOps {
		if method == constants.FieldParameters || method == constants.FieldRef {
			if _, exists := existingOps[method]; !exists {
				existingOps[method] = op
				mergedAny = true
			}
			continue
		}

		if _, methodExists := existingOps[method]; methodExists {
			a.logger.Warn("路径 %s 的方法 %s 在多个服务中重复定义 (当前: %s)，保留首次加载的定义", path, method, serviceName)
		} else {
			existingOps[method] = op
			mergedAny = true
			a.logger.Debug("添加方法 %s 到路径 %s (来自: %s)", method, path, serviceName)
		}
	}

	if !mergedAny {
		a.logger.Debug("路径 %s 的所有方法已存在，无需合并 (来自: %s)", path, serviceName)
	}
}

// mergeDefinitions 合并服务规范中的定义
func (a *Aggregator) mergeDefinitions(specMap map[string]interface{}, serviceName string, allDefinitions map[string]interface{}) {
	definitions, ok := specMap[constants.FieldDefs].(map[string]interface{})
	if !ok {
		return
	}

	for finalDefName, definition := range definitions {
		convertedDef := mathx.ConvertMapKeysToString(definition)
		if existingDef, exists := allDefinitions[finalDefName]; exists {
			a.checkDefinitionConsistency(finalDefName, existingDef, convertedDef, serviceName)
			continue
		}
		allDefinitions[finalDefName] = convertedDef
	}
}

// checkDefinitionConsistency 检查定义一致性
func (a *Aggregator) checkDefinitionConsistency(defName string, existingDef, newDef interface{}, serviceName string) {
	existingJSON, _ := json.Marshal(existingDef)
	newJSON, _ := json.Marshal(newDef)
	if string(existingJSON) != string(newJSON) {
		a.logger.Warn("类型 %s 在不同服务中定义不一致！当前使用第一个定义，忽略 %s 的定义", defName, serviceName)
	} else {
		a.logger.Debug("类型 %s 已存在且定义一致，跳过 (来自: %s)", defName, serviceName)
	}
}

// mergeServiceSpecTags 合并服务规范中的标签
func (a *Aggregator) mergeServiceSpecTags(specMap map[string]interface{}, serviceName string, allTags *[]interface{}, tagNames map[string]bool) {
	tags, ok := specMap[constants.FieldTags].([]interface{})
	if !ok {
		return
	}

	for _, tag := range tags {
		tagMap, ok := tag.(map[string]interface{})
		if !ok {
			continue
		}

		name, exists := tagMap[constants.FieldName]
		if !exists {
			continue
		}

		nameStr := convert.MustString(name)
		if a.addUniqueTag(nameStr, tagMap, allTags, tagNames) {
			a.logger.Debug("添加原始Swagger标签: %s (服务: %s)", nameStr, serviceName)
		}
	}
}

// addUniqueTag 添加唯一标签（通用去重逻辑）
func (a *Aggregator) addUniqueTag(tagKey string, tag interface{}, allTags *[]interface{}, tagSet map[string]bool) bool {
	if tagKey == "" || tagSet[tagKey] {
		return false
	}
	tagSet[tagKey] = true
	*allTags = append(*allTags, tag)
	return true
}

// newMergeError 创建合并错误
func newMergeError(serviceName, reason string) error {
	return errorx.NewError(serrors.ErrTypeAggregateFailed, "服务 "+serviceName+" 合并失败: "+reason)
}
