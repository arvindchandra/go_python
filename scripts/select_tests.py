#!/usr/bin/env python3
import json, sys, yaml
cfg = yaml.safe_load(open(sys.argv[1]))
changed = sys.argv[2:]
services = set(); tests = set()
for prefix, entry in cfg.get('critical_paths', {}).items():
    for f in changed:
        if f.startswith(prefix):
            services.update(entry.get('services', []))
            tests.update(entry.get('tests', []))
print(json.dumps({'services': sorted(services), 'tests': sorted(tests), 'changed': changed}, indent=2))
