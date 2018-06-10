from prometheus_client import start_http_server, Metric, REGISTRY, Gauge
import json
import requests
import sys
import time

class JsonCollector(object):
  def __init__(self, endpoint):
    self._endpoint = endpoint
  def collect(self):
    # Fetch the status
    response = json.loads(requests.get(self._endpoint + "/management/health").content.decode('UTF-8'))
    print response["status"]
    if response['status'] == 'UP':
      metric = Metric('springxd_up', 'SpringXD server status', 'gauge')
      metric.add_sample('springxd_up', value=1.0, labels={})
      yield metric

    # Fetch the metrics   
    response = json.loads(requests.get(self._endpoint + "/management/metrics").content.decode('UTF-8'))

    #collect the metrics
    metric = Metric('springxd_mem_total', 'SpringXD server memory', 'gauge')
    metric.add_sample('springxd_mem_total', value=response["mem"], labels={})
    yield metric
    metric = Metric('springxd_mem_free', 'SpringXD free server memory', 'gauge')
    metric.add_sample('springxd_mem_free', value=response["mem.free"], labels={})
    yield metric
    metric = Metric('springxd_instance_uptime', 'SpringXD instance uptime', 'gauge')
    metric.add_sample('springxd_instance_uptime', value=response["instance.uptime"], labels={})
    yield metric
    metric = Metric('springxd_threads_peak', 'SpringXD Peak Threads', 'gauge')
    metric.add_sample('springxd_threads_peak', value=response["threads.peak"], labels={})
    yield metric
    metric = Metric('springxd_threads_daemon', 'SpringXD Daemon Threads', 'gauge')
    metric.add_sample('springxd_threads_daemon', value=response["threads.daemon"], labels={})
    yield metric
    metric = Metric('springxd_threads_current', 'SpringXD Current Threads', 'gauge')
    metric.add_sample('springxd_threads_peak', value=response["threads"], labels={})
    yield metric
    metric = Metric('springxd_datasource_active', 'SpringXD Datasource Active', 'gauge')
    metric.add_sample('springxd_datasource_active', value=response["datasource.primary.active"], labels={})
    yield metric
    metric = Metric('springxd_datasource_usage', 'SpringXD Datasource Usage', 'gauge')
    metric.add_sample('springxd_datasource_usage', value=response["datasource.primary.usage"], labels={})
    yield metric    
    metric = Metric('springxd_gc',
        'SpringXD GC', 'summary')
    metric.add_sample('springxd_gc_count',
        value=response['gc.g1_young_generation.count'], labels={'gen': 'g1 young'})
    metric.add_sample('springxd_gc_time',
        value=response['gc.g1_young_generation.time'] , labels={'gen': 'g1 young'})
    metric.add_sample('springxd_gc_count',
        value=response['gc.g1_old_generation.count'], labels={'gen': 'g1 old'})
    metric.add_sample('springxd_gc_time',
        value=response['gc.g1_old_generation.time'] , labels={'gen': 'g1 old'})
    yield metric



if __name__ == '__main__':
  # Usage: json_exporter.py port endpoint
  start_http_server(int(sys.argv[1]))
  REGISTRY.register(JsonCollector(sys.argv[2]))

  while True: 
    time.sleep(1)