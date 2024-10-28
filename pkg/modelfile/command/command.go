/*
 *     Copyright 2024 The CNAI Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

// Define the command strings for modelfile.
const (
	// CONFIG is the command to set the configuration of the model, which is used for
	// the model to be served, such as the config.json, generation_config.json, etc.
	// The CONFIG command can be used multiple times in a modelfile, it
	// will be copied the config file to the artifact package as a layer.
	CONFIG = "config"

	// MODEL is the command to set the model file path. The value of this command
	// is the regex of the model file path to match the model file name.
	// The MODEL command can be used multiple times in a modelfile, it will scan
	// the model file path by the regex and copy each model file to the artifact
	// package, and each model file will be a layer.
	MODEL = "model"

	// NAME is the command to set the model name, such as llama3-8b-instruct, gpt2-xl,
	// qwen2-vl-72b-instruct, etc.
	NAME = "name"

	// ARCH is the command to set the architecture of the model, such as transformer,
	// cnn, rnn, etc.
	ARCH = "arch"

	// FAMILY is the command to set the family of the model, such as llama3, gpt2, qwen2, etc.
	FAMILY = "family"

	// FORMAT is the command to set the format of the model, such as onnx, tensorflow, pytorch, etc.
	FORMAT = "format"

	// PARAMSIZE is the command to set the parameter size of the model.
	PARAMSIZE = "paramsize"

	// PRECISION is the command to set the precision of the model, such as bf16, fp16, int8, etc.
	PRECISION = "precision"

	// QUANTIZATION is the command to set the quantization of the model, such as awq, gptq, etc.
	QUANTIZATION = "quantization"
)

// Commands is a list of all the commands that can be used in a modelfile.
var Commands = []string{
	CONFIG,
	MODEL,
	NAME,
	ARCH,
	FAMILY,
	FORMAT,
	PARAMSIZE,
	PRECISION,
	QUANTIZATION,
}