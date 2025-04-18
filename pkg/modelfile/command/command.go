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
	CONFIG = "CONFIG"

	// MODEL is the command to set the model file path. The value of this command
	// is the glob of the model file path to match the model file name.
	// The MODEL command can be used multiple times in a modelfile, it will scan
	// the model file path by the glob and copy each model file to the artifact
	// package, and each model file will be a layer.
	MODEL = "MODEL"

	// CODE is the command to set the code file path. The value of this commands
	// is the glob of the code file path to match the code file name.
	// The CODE command can be used multiple times in a modelfile, it will scan
	// the code file path by the glob and copy each code file to the artifact
	// package, and each code file will be a layer.
	CODE = "CODE"

	// DATASET is the command to set the dataset file path. The value of this commands
	// is the glob of the dataset file path to match the dataset file name.
	// The DATASET command can be used multiple times in a modelfile, it will scan
	// the dataset file path by the glob and copy each dataset file to the artifact
	// package, and each dataset file will be a layer.
	DATASET = "DATASET"

	// DOC is the command to set the documentation file path. The value of this commands
	// is the glob of the documentation file path to match the documentation file name.
	// The DOC command can be used multiple times in a modelfile, it will scan
	// the documentation file path by the glob and copy each documentation file to the artifact
	// package, and each documentation file will be a layer.
	DOC = "DOC"

	// NAME is the command to set the model name, such as llama3-8b-instruct, gpt2-xl,
	// qwen2-vl-72b-instruct, etc.
	NAME = "NAME"

	// ARCH is the command to set the architecture of the model, such as transformer,
	// cnn, rnn, etc.
	ARCH = "ARCH"

	// FAMILY is the command to set the family of the model, such as llama3, gpt2, qwen2, etc.
	FAMILY = "FAMILY"

	// FORMAT is the command to set the format of the model, such as onnx, tensorflow, pytorch, etc.
	FORMAT = "FORMAT"

	// PARAMSIZE is the command to set the parameter size of the model.
	PARAMSIZE = "PARAMSIZE"

	// PRECISION is the command to set the precision of the model, such as bf16, fp16, int8, etc.
	PRECISION = "PRECISION"

	// QUANTIZATION is the command to set the quantization of the model, such as awq, gptq, etc.
	QUANTIZATION = "QUANTIZATION"
)

// Commands is a list of all the commands that can be used in a modelfile.
var Commands = []string{
	CONFIG,
	MODEL,
	CODE,
	DATASET,
	DOC,
	NAME,
	ARCH,
	FAMILY,
	FORMAT,
	PARAMSIZE,
	PRECISION,
	QUANTIZATION,
}
