#!/bin/bash

# Metal Argument buffers allows us to use large sampler arrays on macos
export MVK_CONFIG_USE_METAL_ARGUMENT_BUFFERS=1

# Enable synchronization validation
export VK_LAYER_ENABLES=VK_VALIDATION_FEATURE_ENABLE_SYNCHRONIZATION_VALIDATION_EXT

go run ./game/cmd
