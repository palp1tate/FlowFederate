import consul


class Consul(object):
    def __init__(self, consul_host="127.0.0.1", consul_port=8500):
        self._consul = consul.Consul(host=consul_host, port=consul_port)
        self._service_counters = {}

    def register(self, address, port, service_name, tags, service_id):
        check = consul.Check.tcp(address, port, "10s", "30s", "15s")
        self._consul.agent.service.register(name=service_name, service_id=service_id, tags=tags,
                                            address=address, port=port, check=check)

    def deregister(self, server_id):
        self._consul.agent.service.deregister(service_id=server_id)

    def _fetch_and_format_services(self, name):
        _, nodes = self._consul.health.service(service=name, passing=True)
        if len(nodes) == 0:
            raise Exception('service is empty.')
        services = []
        for node in nodes:
            service = node.get('Service')
            service_address = f"{service['Address']}:{service['Port']}"
            services.append(service_address)
        if not services:
            raise Exception('No available services found for: ' + name)
        return services

    def get_service(self, name):
        services = self._fetch_and_format_services(name)
        if name not in self._service_counters:
            self._service_counters[name] = 0
        service_index = self._service_counters[name] % len(services)
        self._service_counters[name] += 1
        return services[service_index]

    def get_services(self, name, num):
        services = self._fetch_and_format_services(name)
        if num <= 0:
            return []
        selected_services = []
        available_services = list(range(len(services)))
        while len(selected_services) < num and available_services:
            if name not in self._service_counters:
                self._service_counters[name] = 0
            service_index = self._service_counters[name] % len(available_services)
            selected_index = available_services.pop(service_index)
            selected_services.append(services[selected_index])
            self._service_counters[name] += 1
        return selected_services
