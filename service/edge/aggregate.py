import logging
from typing import Iterable

import numpy as np
import torch
from scipy.stats import pearsonr
from torch import nn


def fedavg_aggregate(modules, device):
    num_modules = len(modules)
    avg_state_dict = {}

    # Initialize accumulator dictionary with zeros
    first_module = next(iter(modules))
    for name, param in first_module.state_dict().items():
        avg_state_dict[name] = torch.zeros_like(param).to(device)

    # Sum all model parameters
    for module in modules:
        module_state_dict = module.state_dict()
        for name, param in module_state_dict.items():
            if param.requires_grad:
                avg_state_dict[name] += param.data.to(device)
            else:
                avg_state_dict[name] += param.to(device)

    # Average the parameters
    for name in avg_state_dict:
        if avg_state_dict[name].requires_grad:
            avg_state_dict[name] /= num_modules

    return avg_state_dict


def krum_aggregate(modules, num_modules, device):
    grad_vectors = []

    for module in modules:
        vector = torch.tensor([]).to(device)
        for param in module.parameters():
            if param.requires_grad:
                vector = torch.cat((vector, param.data.view(-1)))
        grad_vectors.append(vector)
    logging.info("Grad vectors loaded successfully.")

    # 计算每个模型与其他所有模型的距离和
    scores = []
    for i, vec1 in enumerate(grad_vectors):
        distances = [
            torch.norm(vec1 - vec2, p=2)
            for j, vec2 in enumerate(grad_vectors)
            if i != j
        ]
        scores.append(sum(sorted(distances)[: num_modules - 2]))
    logging.info("Scores calculated successfully.")

    # 选择得分最低的模型（最接近其他模型的模型）
    selected_index = torch.argmin(torch.tensor(scores))
    selected_state_dict = modules[selected_index].state_dict()
    aggregated_state_dict = {}
    logging.info("Selected model loaded successfully.")

    # Initialize aggregated_state_dict with zeros
    for name, param in selected_state_dict.items():
        aggregated_state_dict[name] = torch.zeros_like(param).to(device)
    logging.info("Aggregated state dict initialized successfully.")

    # Copy relevant parameters from selected model
    for name, param in selected_state_dict.items():
        if param.requires_grad:
            aggregated_state_dict[name].copy_(param.data)
        else:
            aggregated_state_dict[name].copy_(param)
    logging.info("Parameters copied successfully.")
    return aggregated_state_dict


def trimmed_mean_aggregate(modules, device, poisoner_nums):
    # Collect all parameter vectors from all models
    all_params = {name: [] for name, _ in next(iter(modules)).named_parameters()}
    non_update_params = {
        name: None
        for name, param in next(iter(modules)).named_parameters()
        if not param.requires_grad
    }

    for module in modules:
        state_dict = module.state_dict()
        for name, param in state_dict.items():
            if param.requires_grad:
                all_params[name].append(param.view(-1).to(device))
            else:
                non_update_params[name] = param

    logging.info("Parameters flattened and moved to device successfully.")
    trimmed_mean_params = {}

    for name, vectors in all_params.items():
        if len(vectors) == 0:
            # logging.warning(f"No vectors for parameter '{name}', skipping this parameter.")
            continue

        matrix = torch.stack(vectors)
        sorted_matrix, _ = torch.sort(matrix, dim=0)
        trimmed_matrix = (
            sorted_matrix[poisoner_nums:-poisoner_nums]
            if poisoner_nums > 0
            else sorted_matrix
        )

        if trimmed_matrix.size(0) == 0:
            # logging.warning(f"All values trimmed for parameter '{name}', skipping this parameter.")
            continue

        trimmed_mean = trimmed_matrix.mean(dim=0)
        original_shape = next(iter(modules)).state_dict()[name].shape
        trimmed_mean_params[name] = trimmed_mean.view(original_shape)

    # Add non-update parameters back into the result
    for name, param in non_update_params.items():
        if param is not None:
            trimmed_mean_params[name] = param

    logging.info("Trimmed mean computed successfully.")
    return trimmed_mean_params


def median_aggregate(modules, device):
    # Initialize the parameter dictionaries
    all_params = {name: [] for name, _ in next(iter(modules)).named_parameters()}
    non_update_params = {
        name: None
        for name, param in next(iter(modules)).named_parameters()
        if not param.requires_grad
    }
    logging.info("Parameter dictionaries initialized successfully.")
    for module in modules:
        state_dict = module.state_dict()
        for name, param in state_dict.items():
            if param.requires_grad:
                all_params[name].append(param.view(-1).to(device))
            else:
                non_update_params[name] = param
    logging.info("Parameters flattened and moved to device successfully.")
    # Compute the median for each parameter
    median_params = {}

    for name, vectors in all_params.items():
        if len(vectors) == 0:
            # logging.warning(f"No vectors for parameter '{name}', skipping this parameter.")
            continue

        matrix = torch.stack(vectors)
        sorted_matrix, _ = torch.sort(matrix, dim=0)
        median_index = sorted_matrix.size(0) // 2
        if sorted_matrix.size(0) % 2 == 0:
            median = (
                sorted_matrix[median_index - 1] + sorted_matrix[median_index]
            ) / 2.0
        else:
            median = sorted_matrix[median_index]
        original_shape = next(iter(modules)).state_dict()[name].shape
        median_params[name] = median.view(original_shape)
    logging.info("Median computed successfully.")
    # Add non-update parameters from the first model
    for name, param in non_update_params.items():
        if param is not None:
            median_params[name] = param
    logging.info("Non-update parameters added successfully.")
    return median_params


def pefl_aggregate(modules, device):
    # Initialize the parameter dictionaries
    all_params = {name: [] for name, _ in next(iter(modules)).named_parameters()}
    non_update_params = {
        name: None
        for name, param in next(iter(modules)).named_parameters()
        if not param.requires_grad
    }

    for module in modules:
        state_dict = module.state_dict()
        for name, param in state_dict.items():
            if param.requires_grad:
                all_params[name].append(param.view(-1).to(device))
            else:
                non_update_params[name] = param

    # Calculate the median vector
    median_params = {}
    for name, vectors in all_params.items():
        if len(vectors) == 0:
            # logging.warning(f"No vectors for parameter '{name}', skipping this parameter.")
            continue

        matrix = torch.stack(vectors)
        median_index = matrix.size(0) // 2
        sorted_matrix, _ = torch.sort(matrix, dim=0)
        if matrix.size(0) % 2 == 1:
            median = sorted_matrix[median_index]
        else:
            median = (sorted_matrix[median_index - 1] + sorted_matrix[median_index]) / 2
        median_params[name] = median

    # Calculate Pearson correlation and weight updates
    weighted_updates = {
        name: torch.zeros_like(param).to(device)
        for name, param in next(iter(modules)).named_parameters()
        if param.requires_grad
    }
    total_weight = 0
    for module in modules:
        module_state_dict = module.state_dict()
        correlation = 0
        count = 0
        for name, param in module_state_dict.items():
            if name in median_params:
                flat_param = param.view(-1).to(device)
                correlation += pearsonr(
                    flat_param.cpu().numpy(), median_params[name].cpu().numpy()
                )[0]
                count += 1
        correlation /= count if count > 0 else 1
        weight = max(0, np.log((1 + correlation) / (1 - correlation)) - 0.5)
        total_weight += weight
        for name, param in module_state_dict.items():
            if name in median_params:
                weighted_updates[name] += weight * param.view(-1).to(device)

    # Normalize weighted updates
    for name in weighted_updates:
        if total_weight > 0:
            weighted_updates[name] /= total_weight
        # Reshape back to the original shape
        original_shape = next(iter(modules)).state_dict()[name].shape
        weighted_updates[name] = weighted_updates[name].view(original_shape)

    # Add non-update parameters from the first model
    for name, param in non_update_params.items():
        if param is not None:
            weighted_updates[name] = param

    return weighted_updates


def shieldfl_aggregate(
    modules: Iterable[torch.nn.Module], global_model: torch.nn.Module, device: str
):
    cos = nn.CosineSimilarity(dim=-1)
    grad = {}
    grad_vector = {}

    for module in modules:
        diff = {
            name: param - global_model.state_dict()[name].to(device)
            for name, param in module.state_dict().items()
            if param.requires_grad
        }
        if not diff:  # Check if diff is empty
            continue
        vector = torch.cat([data.reshape(-1) for data in diff.values()], dim=0)
        if torch.norm(vector) - 1.0 < 1e-5:
            grad[id(module)] = diff
            grad_vector[id(module)] = vector

    if not grad_vector:  # Check if grad_vector is empty
        logging.warning("No valid gradient vectors found for aggregation.")
        return global_model.state_dict()

    global_grad = torch.zeros_like(next(iter(grad_vector.values())))

    min_cos_similarity = 2
    baseline_grad = torch.zeros_like(global_grad)

    for _id, vec in grad_vector.items():
        cur_cos = cos(vec, global_grad)
        if cur_cos < min_cos_similarity:
            min_cos_similarity = cur_cos
            baseline_grad = vec

    mu = {}
    all_mu = 0
    for _id, vec in grad_vector.items():
        mu[_id] = 1 - cos(baseline_grad, vec)
        all_mu += mu[_id]

    weight_accumulator = {
        name: torch.zeros_like(param, device=device)
        for name, param in global_model.state_dict().items()
    }

    for _id, diff in grad.items():
        for name, data in diff.items():
            weight_accumulator[name].add_((mu[_id] / all_mu) * data)
            if "running_mean" in name or "running_var" in name:
                weight_accumulator[name] /= 1.0  # Normalize with same factor

    # Aggregate weights into the global model
    for name, param in global_model.state_dict().items():
        param.data.copy_(weight_accumulator[name])

    return global_model.state_dict()
