import logging

import grpc
import yaml


def load_configuration(file_path: str) -> dict:
    try:
        with open(file_path, "r") as f:
            return yaml.safe_load(f)
    except Exception as exc:
        raise Exception(f"Failed to load configuration from {file_path}: {exc}")


def load_server_credentials(file_path: str) -> grpc.ServerCredentials:
    try:
        with open(file_path + "/" + "server.crt", "rb") as cf, open(
            file_path + "/" + "server.key", "rb"
        ) as kf, open(file_path + "/" + "ca.crt", "rb") as caf:
            certificate_chain = cf.read()
            private_key = kf.read()
            ca_cert = caf.read()
        return grpc.ssl_server_credentials(
            ((private_key, certificate_chain),),
            root_certificates=ca_cert,
            require_client_auth=True,
        )
    except Exception as exc:
        logging.error(f"Failed to load credentials from {file_path}: {exc}")
        raise


def load_client_credentials(file_path: str) -> grpc.ChannelCredentials:
    """Load SSL credentials for gRPC client."""
    with open(file_path + "/" + "server.crt", "rb") as scf, open(
        file_path + "/" + "client.crt", "rb"
    ) as ccf, open(file_path + "/" + "client.key", "rb") as kf:
        trusted_certs = scf.read()
        client_cert = ccf.read()
        client_key = kf.read()
    return grpc.ssl_channel_credentials(
        root_certificates=trusted_certs,
        private_key=client_key,
        certificate_chain=client_cert,
    )
