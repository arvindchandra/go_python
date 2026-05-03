#!/usr/bin/env python3
import json
import sys
import xml.etree.ElementTree as ET

path = sys.argv[1]
root = ET.parse(path).getroot()
summary = {'tests': 0, 'failures': 0, 'errors': 0, 'skipped': 0}
if root.tag == 'testsuite':
    summary['tests'] = int(root.attrib.get('tests', 0))
    summary['failures'] = int(root.attrib.get('failures', 0))
    summary['errors'] = int(root.attrib.get('errors', 0))
    summary['skipped'] = int(root.attrib.get('skipped', 0))
else:
    for suite in root.findall('testsuite'):
        summary['tests'] += int(suite.attrib.get('tests', 0))
        summary['failures'] += int(suite.attrib.get('failures', 0))
        summary['errors'] += int(suite.attrib.get('errors', 0))
        summary['skipped'] += int(suite.attrib.get('skipped', 0))
print(json.dumps(summary, indent=2))
