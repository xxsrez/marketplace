import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
TASK_MANAGER_SKILL = (
    REPOSITORY_ROOT
    / "plugins"
    / "task-manager"
    / "skills"
    / "task-manager"
    / "SKILL.md"
)


class TaskManagerFileBridgePackagingTest(unittest.TestCase):
    def test_skill_prefers_the_hosted_codex_file_bridge(self) -> None:
        skill = TASK_MANAGER_SKILL.read_text(encoding="utf-8")
        normalized = " ".join(skill.split())

        self.assertIn("connector-first route", normalized)
        self.assertIn("host-facing absolute path", normalized)
        self.assertIn("remote MCP never receives or reads that path", normalized)
        self.assertIn("description or native Comment body", normalized)
        self.assertNotIn(
            "absolute filesystem path can go only to the local",
            normalized,
        )


if __name__ == "__main__":
    unittest.main()
