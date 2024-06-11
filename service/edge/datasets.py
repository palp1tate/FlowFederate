import json
import time

from torchvision import datasets, transforms

import numpy as np
import matplotlib.pyplot as plt
from torchvision import datasets

seed = 42


def get_dataset(dir, name):
    if name == "mnist":
        train_dataset = datasets.MNIST(
            dir, train=True, download=True, transform=transforms.ToTensor()
        )
        eval_dataset = datasets.MNIST(dir, train=False, transform=transforms.ToTensor())

    if name == "fashionmnist":
        train_dataset = datasets.FashionMNIST(
            dir, train=True, download=True, transform=transforms.ToTensor()
        )
        eval_dataset = datasets.FashionMNIST(
            dir, train=False, transform=transforms.ToTensor()
        )
        print("get_fashionmnist")

    elif name == "cifar10":
        transform_train = transforms.Compose(
            [
                transforms.RandomCrop(32, padding=4),
                transforms.RandomHorizontalFlip(),
                transforms.ToTensor(),
                transforms.Normalize(
                    (0.4914, 0.4822, 0.4465), (0.2023, 0.1994, 0.2010)
                ),
            ]
        )

        transform_test = transforms.Compose(
            [
                transforms.ToTensor(),
                transforms.Normalize(
                    (0.4914, 0.4822, 0.4465), (0.2023, 0.1994, 0.2010)
                ),
            ]
        )

        train_dataset = datasets.CIFAR10(
            dir, train=True, download=True, transform=transform_train
        )
        eval_dataset = datasets.CIFAR10(dir, train=False, transform=transform_test)

    elif name == "cifar100":
        transform_train = transforms.Compose(
            [
                transforms.RandomCrop(32, padding=4),
                transforms.RandomHorizontalFlip(),
                transforms.ToTensor(),
                transforms.Normalize(
                    (0.4914, 0.4822, 0.4465), (0.2023, 0.1994, 0.2010)
                ),
            ]
        )

        transform_test = transforms.Compose(
            [
                transforms.ToTensor(),
                transforms.Normalize(
                    (0.4914, 0.4822, 0.4465), (0.2023, 0.1994, 0.2010)
                ),
            ]
        )

        train_dataset = datasets.CIFAR100(
            dir, train=True, download=True, transform=transform_train
        )
        eval_dataset = datasets.CIFAR100(dir, train=False, transform=transform_test)

    return train_dataset, eval_dataset


# *生成一组非独立同分布（non-IID）的数据集，其中每个客户端的数据集包含了多个类别的样本。
# *深度学习中，非独立同分布（non-IID）数据集是指数据集中的样本分布不同于整个数据集的分布。如果数据集是非独立同分布的，模型就可能无法很好地学习到数据集的分布特征，从而影响模型的性能和泛化能力。
# *在实际应用中，由于数据来源的多样性和数据采集的不确定性，很多数据集都是非独立同分布的。在联合学习等场景中，需要将非独立同分布的数据集分配给不同的客户端进行训练，以提高模型的泛化能力和性能。
def get_nonIID_data(conf):
    # *用于存储每个客户端的数据集合
    client_idx = {}
    # *包含所有类别的列表all_data，其中每个元素代表一个类别
    all_data = []
    for i in range(conf["classes"]):
        all_data.append(i)
    for i in range(conf["clients"]):
        # *从 all_data 中随机选择 conf["client_classes"] 个类别，生成了一个包含了多个类别的样本集合 samples
        # *将 samples 添加到 client_idx 中，生成了一个多个类别的样本集合。最终，代码返回了 client_idx，其中包含了所有客户端的数据集合。
        # *这里使用replace=False参数来确保每个类别只被选择一次，避免重复
        samples = np.random.choice(all_data, size=conf["client_classes"], replace=False)
        client_idx[i + 1] = samples

    # *返回client_idx字典，其中包含了所有客户端的数据集合。每个客户端的数据集合是一个包含多个类别的样本集合，用于构建非独立同分布的数据集。
    # *也就是说这个函数完成的任务是：传入配置信息，执行非独立同分布数据集合，每个客户端都有随机选取的配置信息指定数目的size=conf["client_classes"]个data
    return client_idx


def dirichlet_split_noniid(train_labels, alpha, n_clients):
    n_classes = train_labels.max() + 1
    # (K, N) 类别标签分布矩阵X，记录每个类别划分到每个client去的比例
    # *使用Dirichlet分布生成标签分布，以便在联邦学习中模拟非独立同分布（Non-IID）数据
    # *np.random.dirichlet()是NumPy库中的函数，用于从Dirichlet分布中抽取样本。在这里，np.random.dirichlet([alpha] * n_clients, n_classes)
    # *生成一个二维形状为(n_clients, n_classes)的NumPy数组，每个元素表示每个客户端在每个类别上的标签分布。每一行表示一个客户端的标签分布，每一列表示一个类别的概率，
    # *alpha是Dirichlet分布的参数，用于控制标签分布的平滑度。
    # *Dirichlet分布是一种连续多元概率分布，通常用于对未知概率向量的不确定性建模，Dirichlet分布样本是一个概率向量，其各个分量之和为1
    label_distribution = np.random.dirichlet([alpha] * n_clients, n_classes)
    # label_distribution = np.random.dirichlet(np.repeat(alpha, n_clients))
    # (K, ...) 记录K个类别对应的样本索引集合
    # *将训练集中每个类别的标签索引存储在一个列表中，列表中的第i个元素是一个包含训练集中标签为i的样本索引的NumPy数组。np.argwhere()返回满足条件的元素的索引
    # *np.argwhere(train_labels == y)返回一个形状为(m, 1)的NumPy数组，其中m是训练集中标签为y的样本数量，每个元素是一个标签为y的样本的索引,然后展平为一维数组

    """
    例如，如果训练标签中有100个样本，其中有30个属于类别y，那么类别y的样本索引数组就是一个包含30个元素的一维数组，表示这30个样本在训练数据集中的索引位置。
    """
    class_idcs = [np.argwhere(train_labels == y).flatten() for y in range(n_classes)]
    # 记录N个client分别对应的样本索引集合
    # *使用zip(class_idcs, label_distribution)将标签索引和标签分布进行迭代。在每次迭代中，k_idcs表示类别k的样本索引数组，fracs表示标签k在每个客户端上的样本比例。
    # *client_idcs列表被初始化为包含了n_clients个空列表的列表。每个空列表都用于存储一个客户端对应的样本索引集合。
    client_idcs = [[] for _ in range(n_clients)]
    for k_idcs, fracs in zip(class_idcs, label_distribution):
        # np.split按照比例fracs将标签为k的样本索引k_idcs划分为了N个子集
        # i表示第i个client，idcs表示其对应的样本索引集合idcs
        # *将标签k的样本索引数组按照比例fracs划分为多个子集。这里使用np.cumsum(fracs)[:-1] * len(k_idcs)计算每个划分点的位置，然后使用np.split()函数进行划分。
        # *划分后的每个子集都是一个客户端对应的样本索引数组。
        # *`enumerate()`函数同时获取索引`i`和划分后的子集样本索引数组`idcs`，然后将`idcs`添加到`client_idcs`列表中。
        for i, idcs in enumerate(
            np.split(k_idcs, (np.cumsum(fracs)[:-1] * len(k_idcs)).astype(int))
        ):
            # *每个客户端对应的样本索引数组添加到client_idcs列表中。
            client_idcs[i] += [idcs]
    # *将每个客户端的样本索引数组连接起来，得到一个包含每个客户端样本索引的列表
    # *通过列表推导式[np.concatenate(idcs) for idcs in client_idcs]，遍历client_idcs列表中的每个元素idcs，
    # *使用np.concatenate()函数将其中的样本索引数组连接起来。
    # *最终，得到的client_idcs列表中的每个元素都是一个包含了一个客户端对应的样本索引的一维数组。
    # *这个操作之前client_idcs是一个二维列表，其中每个元素都是一个包含了一个客户端对应的样本索引的一维数组。
    client_idcs = [np.concatenate(idcs) for idcs in client_idcs]
    # return client_idcs
    # *每个客户端的样本索引数组存储到一个字典client_idx中，其中键是客户端的编号，值是客户端对应的样本索引数组。
    client_idx = {}
    for i in range(len(client_idcs)):
        client_idx[i + 1] = client_idcs[i]
    return client_idx


def dirichlet_nonIID_data(train_data, conf):
    # *NumPy库中的函数，用于设置随机数生成器的种子(seed)值。其中，seed是一个整数，用于指定随机数生成器的种子值。
    # *这个种子值可以使随机数生成器的输出可预测，因为每次使用相同的种子值，都会得到相同的随机数序列
    np.random.seed(seed)

    # *得到传入参数——训练数据集的类别
    classes = train_data.classes
    # *获得类别总数
    n_classes = len(classes)
    # *注释代码作用是：将PyTorch数据集中的训练集和测试集的标签合并为一个NumPy数组，标签沿着0轴连接
    # labels = np.concatenate([np.array(train_data.targets), np.array(test_data.targets)], axis=0)
    # *数据集中的标签转换为NumPy数组
    labels = np.array(train_data.targets)
    # dataset = ConcatDataset([train_data, test_data])

    # 我们让每个client不同label的样本数量不同，以此做到Non-IID划分
    # client_idcs = dirichlet_split_noniid(labels, alpha=conf["dirichlet_alpha"], n_clients=conf["clients"])
    return dirichlet_split_noniid(
        labels, alpha=conf["dirichlet_alpha"], n_clients=conf["clients"]
    )

    # 展示不同label划分到不同client的情况
    fig = plt.figure(figsize=(8, 7))
    plt.hist(
        [labels[idc] for idc in client_idcs],
        stacked=True,
        bins=np.arange(min(labels) - 0.5, max(labels) + 1.5, 1),
        label=["Terminal {}".format(i) for i in range(conf["clients"])],
        rwidth=0.8,
    )
    # plt.xticks(np.arange(n_classes), fontsize=20)
    plt.xticks(np.array([0, 9, 19, 29, 39, 49, 59, 69, 79, 89, 99]), fontsize=20)
    plt.xlabel("Label type", fontsize=20)
    plt.ylabel("Number of samples", fontsize=20)
    plt.legend(loc="upper right")
    plt.title("CIFAR100", fontsize=20)
    plt.show()

    filename = "Lai-CIFAR100.pdf"
    fig.savefig("../result/" + filename, bbox_inches="tight")


if __name__ == "__main__":
    with open("../utils/conf.json", "r") as f:
        conf = json.load(f)

    train_datasets, eval_datasets = get_dataset("../data/", "cifar100")

    dirichlet_nonIID_data(train_datasets, conf)
