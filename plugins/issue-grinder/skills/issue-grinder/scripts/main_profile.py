#!/usr/bin/env python3
"""Read only this Codex session's effective profile; never infer it from config."""
from __future__ import annotations
import json
import os
from pathlib import Path
import re


def read_profile(codex_home: Path, thread_id: str) -> dict:
    if not re.fullmatch(r'[0-9a-fA-F]{8}(?:-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}', thread_id):
        return {'status': 'unknown', 'reason': 'current_thread_id_unavailable'}
    candidates = []
    for folder in ('sessions', 'archived_sessions'):
        base = codex_home / folder
        if base.is_dir():
            candidates.extend(base.rglob(f'*-{thread_id}.jsonl'))
    for path in sorted(candidates, key=lambda p: p.stat().st_mtime, reverse=True):
        try:
            with path.open(encoding='utf-8') as handle:
                meta = json.loads(handle.readline())
                if meta.get('type') != 'session_meta' or meta.get('payload', {}).get('id') != thread_id:
                    continue
                latest = None
                for line in handle:
                    try:
                        record = json.loads(line)
                    except json.JSONDecodeError:
                        continue  # a concurrently appended final line may be incomplete
                    if record.get('type') == 'turn_context':
                        latest = record.get('payload') or {}
                if latest and latest.get('model') and latest.get('effort'):
                    return {'status': 'observed', 'thread_id': thread_id,
                            'model': latest['model'], 'effort': latest['effort'],
                            'source': 'current_session_latest_turn_context'}
        except (OSError, json.JSONDecodeError):
            continue
    return {'status': 'unknown', 'thread_id': thread_id, 'reason': 'current_profile_not_observed'}


def main() -> int:
    result = read_profile(Path(os.environ.get('CODEX_HOME') or Path.home()/'.codex'),
                          os.environ.get('CODEX_THREAD_ID', ''))
    print(json.dumps(result, sort_keys=True))
    return 0 if result['status'] == 'observed' else 2


if __name__ == '__main__':
    raise SystemExit(main())
