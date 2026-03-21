/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-03-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-20 00:00:00
 * @FilePath: \go-swagger\loader\file.go
 * @Description: 文件规范加载器 - 从本地文件加载 Swagger 规范，自动检测 JSON/YAML 格式
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package loader

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/kamalyes/go-swagger/constants"
	serrors "github.com/kamalyes/go-swagger/errors"
	"github.com/kamalyes/go-toolbox/pkg/errorx"
	"gopkg.in/yaml.v3"
)

// LoadSpecFromFile 从文件加载 Swagger 规范（支持 JSON 和 YAML 格式）
func LoadSpecFromFile(filePath string) (map[string]interface{}, error) {
	if !filepath.IsAbs(filePath) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return nil, errorx.NewError(serrors.ErrTypeLoaderReadFailed, filePath)
		}
		filePath = absPath
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errorx.NewError(serrors.ErrTypeLoaderFileNotFound, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errorx.NewError(serrors.ErrTypeLoaderReadFailed, filePath)
	}

	var spec map[string]interface{}
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case constants.FileExtYAML, constants.FileExtYML:
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return nil, errorx.NewError(serrors.ErrTypeLoaderParseFailed, filePath)
		}
	case constants.FileExtJSON:
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, errorx.NewError(serrors.ErrTypeLoaderParseFailed, filePath)
		}
	default:
		return nil, errorx.NewError(serrors.ErrTypeInvalidFileFormat, ext)
	}

	return spec, nil
}

// LoadSpecFromPath 从指定路径加载规范文件，根据扩展名自动检测格式
// 与 LoadSpecFromFile 不同，此方法支持无扩展名时的自动检测回退
func LoadSpecFromPath(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.NewError(serrors.ErrTypeLoaderReadFailed, path)
	}

	var spec map[string]interface{}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case constants.FileExtYAML, constants.FileExtYML:
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return nil, errorx.NewError(serrors.ErrTypeLoaderParseFailed, path)
		}
	case constants.FileExtJSON:
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, errorx.NewError(serrors.ErrTypeLoaderParseFailed, path)
		}
	default:
		if err := json.Unmarshal(data, &spec); err != nil {
			if yamlErr := yaml.Unmarshal(data, &spec); yamlErr != nil {
				return nil, errorx.NewError(serrors.ErrTypeLoaderParseFailed, path)
			}
		}
	}

	return spec, nil
}
