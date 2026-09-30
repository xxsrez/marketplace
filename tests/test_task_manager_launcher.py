import os
import json
import platform
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


class TaskManagerLauncherTest(unittest.TestCase):
    def test_packaged_native_binary_discovery_without_credentials(self):
        if platform.system() not in ("Darwin", "Linux"):
            self.skipTest("unsupported platform")
        requests = [
            {"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-11-25"}},
            {"jsonrpc": "2.0", "id": 2, "method": "tools/list"},
        ]
        if platform.system() == "Linux":
            requests.append({"jsonrpc": "2.0", "id": 3, "method": "tools/call", "params": {
                "name": "attach_local_file_to_task", "arguments": {
                    "taskRef": "TM-1", "fileRef": "fixture", "idempotencyKey": "fixture",
                },
            }})
        env = {key: value for key, value in os.environ.items() if key != "TASK_MANAGER_LOCAL_TOKEN"}
        result = subprocess.run(
            [str(ROOT / "plugins/task-manager/bin/task-manager-local-launcher")],
            input="".join(json.dumps(request) + "\n" for request in requests),
            capture_output=True, text=True, env=env, timeout=10,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        responses = {item["id"]: item for item in map(json.loads, result.stdout.splitlines())}
        self.assertEqual(responses[1]["result"]["serverInfo"]["version"], "0.3.0")
        self.assertEqual({tool["name"] for tool in responses[2]["result"]["tools"]}, {
            "upload_local_file", "attach_local_file_to_task",
        })
        if platform.system() == "Linux":
            self.assertTrue(responses[3]["result"]["isError"])
            self.assertIn("environment_token_required", responses[3]["result"]["content"][0]["text"])

    def test_dispatches_by_os_and_architecture(self):
        with tempfile.TemporaryDirectory(prefix="tm launcher ") as directory:
            root = Path(directory)
            launcher = root / "task-manager-local-launcher"
            shutil.copy2(ROOT / "plugins/task-manager/bin/task-manager-local-launcher", launcher)
            uname = root / "uname"
            uname.write_text('#!/bin/sh\ncase "$1" in -s) echo "$TEST_OS";; -m) echo "$TEST_ARCH";; esac\n')
            uname.chmod(0o755)
            for platform in ("darwin-arm64", "darwin-amd64", "linux-arm64", "linux-amd64"):
                binary = root / f"task-manager-local-{platform}"
                binary.write_text(f'#!/bin/sh\nprintf "%s\\n" "{platform}" "$@"\n')
                binary.chmod(0o755)
            for system, arch, expected in (
                ("Darwin", "arm64", "darwin-arm64"),
                ("Darwin", "x86_64", "darwin-amd64"),
                ("Linux", "x86_64", "linux-amd64"),
                ("Linux", "aarch64", "linux-arm64"),
                ("Linux", "arm64", "linux-arm64"),
            ):
                with self.subTest(system=system, arch=arch):
                    result = subprocess.run(
                        [str(launcher), "argument with spaces"], capture_output=True, text=True,
                        env={**os.environ, "PATH": f"{root}:/usr/bin:/bin", "TEST_OS": system, "TEST_ARCH": arch},
                        timeout=5,
                    )
                    self.assertEqual(result.returncode, 0, result.stderr)
                    self.assertEqual(result.stdout.splitlines(), [expected, "argument with spaces"])
            result = subprocess.run(
                [str(launcher)], capture_output=True, text=True,
                env={**os.environ, "PATH": f"{root}:/usr/bin:/bin", "TEST_OS": "Linux", "TEST_ARCH": "riscv64"},
                timeout=5,
            )
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("Linux riscv64", result.stderr)


if __name__ == "__main__":
    unittest.main()
