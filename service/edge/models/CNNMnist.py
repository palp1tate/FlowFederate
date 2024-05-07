import torch.nn.functional as F
from torch import nn


class cnnmnist(nn.Module):
    def __init__(self):
        super(cnnmnist, self).__init__()
        self.conv1 = nn.Conv2d(1, 10, kernel_size=5)
        self.conv2 = nn.Conv2d(10, 20, kernel_size=5)
        self.conv2_drop = nn.Dropout2d()
        self.fc1 = nn.Linear(320, 50)
        self.fc = nn.Linear(50, 10)

    def forward(self, x):
        x1 = F.relu(F.max_pool2d(self.conv1(x), 2))
        x2 = F.relu(F.max_pool2d(self.conv2_drop(self.conv2(x1)), 2))
        # 数据摊平
        x2 = x2.view(-1, x2.shape[1] * x2.shape[2] * x2.shape[3])
        x3 = F.relu(self.fc1(x2))
        x3 = F.dropout(x3, training=self.training)
        x4 = self.fc(x3)
        return [x1, x2, x3, x4], F.log_softmax(x4, dim=1)


class LeNet(nn.Module):
    def __init__(self):
        super(LeNet, self).__init__()
        self.conv1 = nn.Conv2d(1, 6, 5)
        self.pool1 = nn.MaxPool2d(2, 2)
        self.conv2 = nn.Conv2d(6, 16, 5)
        self.pool2 = nn.MaxPool2d(2, 2)
        self.fc1 = nn.Linear(16 * 4 * 4, 120)
        self.fc2 = nn.Linear(120, 84)
        self.fc = nn.Linear(84, 10)

    def forward(self, x):
        x1 = F.relu(self.conv1(x))
        x2 = self.pool1(x1)
        x3 = F.relu(self.conv2(x2))
        x4 = self.pool2(x3)
        x4 = x4.view(-1, 16 * 4 * 4)
        x5 = F.relu(self.fc1(x4))
        x6 = F.relu(self.fc2(x5))
        x7 = self.fc(x6)
        return [x1, x2, x3, x4, x5, x6, x7], x7


class SimpleCNN(nn.Module):
    def __init__(self):
        super(SimpleCNN, self).__init__()
        self.conv1 = nn.Conv2d(1, 6, 5)  # 调整为 1 个输入通道
        self.relu = nn.ReLU()
        self.pool = nn.MaxPool2d(2, 2)
        self.conv2 = nn.Conv2d(6, 16, 5)

        self.fc1 = nn.Linear(16 * 4 * 4, 120)
        self.fc2 = nn.Linear(120, 84)
        self.fc3 = nn.Linear(84, 10)

    def forward(self, x):
        x1 = self.pool(F.relu(self.conv1(x)))
        x2 = self.pool(F.relu(self.conv2(x1)))
        x2 = x2.view(-1, 16 * 4 * 4)  # 调整 flatten 的尺寸
        x3 = F.relu(self.fc1(x2))
        x4 = F.relu(self.fc2(x3))
        x5 = self.fc3(x4)
        return [x1, x2, x3, x4, x5], x5
