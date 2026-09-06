#!/usr/bin/env python3
"""Fail-closed preflight for Issue Grinder child-agent model routing."""

from __future__ import annotations

import argparse
import hashlib
import json
from dataclasses import asdict, dataclass


SCHEMA = "issue-grinder/model-routing/v2"
LUNA_MODEL = "gpt-5.6-luna"
LUNA_MAX_EFFORT = "max"
LUNA_MODES = frozenset({"economical", "balance"})
BALANCE_SOL_ROLES = frozenset({"specialist", "final_review"})


@dataclass(frozen=True)
class RoutingReceipt:
    schema: str
    packet_id: str
    mode: str
    semantic_role: str
    agent_type: str
    requested_model: str
    requested_effort: str
    fork_turns: str
    dispatch_fingerprint: str
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
    """Return whether the selected mode requires an explicit Luna child."""

    if mode == "balance" and semantic_role in BALANCE_SOL_ROLES:
        return False
    if mode not in LUNA_MODES or user_profile_override:
        return False
    return True


def required_luna_effort(mode: str) -> str:
    """Return the mode-specific default effort for a required Luna child."""

    return LUNA_MAX_EFFORT


def validate_route(
    *,
    packet_id: str,
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

    normalized_packet_id = packet_id.strip()
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
    if not normalized_packet_id:
        defects.append("blank_packet_id")
    if normalized_mode not in {"solo", "classic", *LUNA_MODES}:
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
    if luna_required:
        required_effort = required_luna_effort(normalized_mode)
        if normalized_model != LUNA_MODEL:
            defects.append("luna_model_required")
        if normalized_effort != required_effort:
            defects.append(f"luna_{required_effort}_effort_required")
        if normalized_actual_model is not None and normalized_actual_model != LUNA_MODEL:
            defects.append("actual_luna_model_mismatch")
        if normalized_actual_effort is not None and normalized_actual_effort != required_effort:
            defects.append("actual_luna_effort_mismatch")
    elif normalized_mode == "balance" and normalized_role in BALANCE_SOL_ROLES and not user_profile_override:
        if normalized_model != "gpt-5.6-sol" or normalized_effort != "xhigh":
            defects.append("balance_sol_xhigh_required")
        if normalized_actual_model is not None and (normalized_actual_model != normalized_model or normalized_actual_effort != normalized_effort):
            defects.append("actual_balance_profile_mismatch")
    elif normalized_actual_model is not None:
        if normalized_actual_model != normalized_model:
            defects.append("actual_model_mismatch")
        if normalized_actual_effort != normalized_effort:
            defects.append("actual_effort_mismatch")

    fingerprint_payload = {
        "packet_id": normalized_packet_id,
        "mode": normalized_mode,
        "semantic_role": normalized_role,
        "agent_type": normalized_agent_type,
        "requested_model": normalized_model,
        "requested_effort": normalized_effort,
        "fork_turns": normalized_fork,
        "user_profile_override": user_profile_override,
    }
    dispatch_fingerprint = hashlib.sha256(
        json.dumps(
            fingerprint_payload,
            ensure_ascii=True,
            sort_keys=True,
            separators=(",", ":"),
        ).encode("utf-8")
    ).hexdigest()

    return RoutingReceipt(
        schema=SCHEMA,
        packet_id=normalized_packet_id,
        mode=normalized_mode,
        semantic_role=normalized_role,
        agent_type=normalized_agent_type,
        requested_model=normalized_model,
        requested_effort=normalized_effort,
        fork_turns=normalized_fork,
        dispatch_fingerprint=dispatch_fingerprint,
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
    parser.add_argument("--packet-id", required=True)
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
        packet_id=args.packet_id,
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
