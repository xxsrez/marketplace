#!/usr/bin/env python3
"""Fail-closed preflight for Issue Grinder child-agent model routing."""

from __future__ import annotations

import argparse
import json
from dataclasses import asdict, dataclass


SCHEMA = "issue-grinder/model-routing/v1"
LUNA_MODEL = "gpt-5.6-luna"
LUNA_EFFORT = "max"
ECONOMICAL_MODES = frozenset({"balance", "swarm", "economical"})
FORCED_PROFILE_AGENT_TYPES = frozenset({"critic", "reviewer"})
BALANCE_CONTROLLER_ROLES = frozenset(
    {"material_judgment", "integration_decision", "final_review"}
)


@dataclass(frozen=True)
class RoutingReceipt:
    schema: str
    mode: str
    semantic_role: str
    agent_type: str
    requested_model: str
    requested_effort: str
    fork_turns: str
    user_profile_override: bool
    luna_required: bool
    telemetry_status: str
    actual_model: str | None
    actual_effort: str | None
    allowed: bool
    defects: tuple[str, ...]


def _valid_bounded_fork(value: str) -> bool:
    if value == "none":
        return True
    return value.isdigit() and int(value) > 0


def luna_required_for(
    mode: str,
    semantic_role: str,
    *,
    user_profile_override: bool,
) -> bool:
    """Return whether the selected mode requires an explicit Luna Max child."""

    if mode not in ECONOMICAL_MODES or user_profile_override:
        return False
    if mode == "balance" and semantic_role in BALANCE_CONTROLLER_ROLES:
        return False
    return True


def validate_route(
    *,
    mode: str,
    semantic_role: str,
    agent_type: str,
    model: str,
    effort: str,
    fork_turns: str,
    user_profile_override: bool = False,
    actual_model: str | None = None,
    actual_effort: str | None = None,
) -> RoutingReceipt:
    """Validate requested and, when known, observed child-agent routing."""

    normalized_mode = mode.strip().casefold()
    normalized_role = semantic_role.strip().casefold()
    normalized_agent_type = agent_type.strip().casefold()
    normalized_model = model.strip().casefold()
    normalized_effort = effort.strip().casefold()
    normalized_fork = fork_turns.strip().casefold()
    normalized_actual_model = (
        actual_model.strip().casefold() if actual_model is not None else None
    )
    normalized_actual_effort = (
        actual_effort.strip().casefold() if actual_effort is not None else None
    )

    defects: list[str] = []
    if normalized_mode not in {"solo", "classic", *ECONOMICAL_MODES}:
        defects.append("unknown_mode")
    if not normalized_role:
        defects.append("blank_semantic_role")
    if not normalized_agent_type:
        defects.append("blank_agent_type")
    if not _valid_bounded_fork(normalized_fork):
        defects.append("unbounded_or_missing_fork_turns")
    if (actual_model is None) != (actual_effort is None):
        defects.append("incomplete_actual_profile")

    luna_required = luna_required_for(
        normalized_mode,
        normalized_role,
        user_profile_override=user_profile_override,
    )
    if (
        normalized_mode in ECONOMICAL_MODES
        and normalized_agent_type in FORCED_PROFILE_AGENT_TYPES
        and not user_profile_override
    ):
        defects.append("platform_agent_type_bypasses_mode_profile")
    if luna_required:
        if normalized_model != LUNA_MODEL:
            defects.append("luna_model_required")
        if normalized_effort != LUNA_EFFORT:
            defects.append("luna_max_effort_required")
        if normalized_actual_model is not None and normalized_actual_model != LUNA_MODEL:
            defects.append("actual_luna_model_mismatch")
        if normalized_actual_effort is not None and normalized_actual_effort != LUNA_EFFORT:
            defects.append("actual_luna_effort_mismatch")
    elif normalized_actual_model is not None:
        if normalized_actual_model != normalized_model:
            defects.append("actual_model_mismatch")
        if normalized_actual_effort != normalized_effort:
            defects.append("actual_effort_mismatch")

    return RoutingReceipt(
        schema=SCHEMA,
        mode=normalized_mode,
        semantic_role=normalized_role,
        agent_type=normalized_agent_type,
        requested_model=normalized_model,
        requested_effort=normalized_effort,
        fork_turns=normalized_fork,
        user_profile_override=user_profile_override,
        luna_required=luna_required,
        telemetry_status=(
            "observed"
            if normalized_actual_model is not None
            and normalized_actual_effort is not None
            else "telemetry_pending"
        ),
        actual_model=normalized_actual_model,
        actual_effort=normalized_actual_effort,
        allowed=not defects,
        defects=tuple(defects),
    )


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--mode", required=True)
    parser.add_argument("--semantic-role", required=True)
    parser.add_argument("--agent-type", required=True)
    parser.add_argument("--model", required=True)
    parser.add_argument("--effort", required=True)
    parser.add_argument("--fork-turns", required=True)
    parser.add_argument("--user-profile-override", action="store_true")
    parser.add_argument("--actual-model")
    parser.add_argument("--actual-effort")
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    receipt = validate_route(
        mode=args.mode,
        semantic_role=args.semantic_role,
        agent_type=args.agent_type,
        model=args.model,
        effort=args.effort,
        fork_turns=args.fork_turns,
        user_profile_override=args.user_profile_override,
        actual_model=args.actual_model,
        actual_effort=args.actual_effort,
    )
    print(json.dumps(asdict(receipt), ensure_ascii=False, sort_keys=True))
    return 0 if receipt.allowed else 2


if __name__ == "__main__":
    raise SystemExit(main())
