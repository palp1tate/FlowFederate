import torch

from models.CNNMnist import cnnmnist, LeNet, SimpleCNN
from models.ResNetv1 import resnet18, resnet34, resnet50, resnet101, resnet152
from models.ResNetv2 import resnet8, resnet14, resnet20, resnet32, resnet44, resnet56, \
    resnet110, resnet116, resnet8x4, resnet32x4


# 标准 ResNet 架构（适用于 ImageNet 及其他数据集）：
# resnet18
# resnet34
# resnet50
# resnet101
# resnet152
# 这些架构是基于 ImageNet 数据集设计的，通常用于处理较大的图像（224x224 或更大）。它们在处理 CIFAR-10 和 CIFAR-100 这样的小图像数据集（32x32）时也可以有效工作，但需要适当调整。
# CIFAR ResNet 架构（专为 CIFAR 数据集设计）：
# resnet8
# resnet14
# resnet20
# resnet32
# resnet44
# resnet56
# resnet110
# resnet116
# resnet8x4
# resnet32x4
# 这些架构是专门为 CIFAR-10 和 CIFAR-100 设计的。它们对小尺寸图像（32x32）的输入进行优化，避免了对输入图像尺寸进行额外处理的麻烦。

# 标准 ResNet 架构：
# 适用数据集：ImageNet、CIFAR-10、CIFAR-100 等
# 特点：需要额外的预处理来调整图像尺寸。
# CIFAR ResNet 架构：
# 适用数据集：CIFAR-10、CIFAR-100
# 特点：针对 32x32 输入进行了优化。

# 这里我们是
# resnet18
# resnet34
# resnet50
# resnet101
# resnet152
# 训练cifar10
# 然而
# resnet8
# resnet14
# resnet20
# resnet32
# resnet44
# resnet56
# resnet110
# resnet116
# resnet8x4
# resnet32x4
# 用于cifar100



def get_model(name):
    if name == "resnet18":
        model = resnet18()
    elif name == "resnet34":
        model = resnet34()
    elif name == "resnet50":
        model = resnet50()
    elif name == "resnet101":
        model = resnet101()
    elif name == "resnet152":
        model = resnet152()



    elif name == "resnet8":
        model = resnet8()
    elif name == "resnet14":
        model = resnet14()
    elif name == "resnet20":
        model = resnet20()
    elif name == "resnet32":
        model = resnet32()
    elif name == "resnet44":
        model = resnet44()
    elif name == "resnet56":
        model = resnet56()
    elif name == "resnet110":
        model = resnet110()
    elif name == "resnet116":
        model = resnet116()
    elif name == "resnet8x4":
        model = resnet8x4()
    elif name == "resnet32x4":
        model = resnet32x4()



    elif name == "cnnmnist":
        model = cnnmnist()
    elif name == "lenet":
        model = LeNet()
    elif name == "simple-cnn":
        model = SimpleCNN()

    if torch.cuda.is_available():
        return model.cuda()
    else:
        return model
