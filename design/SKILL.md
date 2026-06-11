# SKILL — moved

The Atempo design skill is now wired into Claude Code as a real project skill:

**`.claude/skills/atempo-design/SKILL.md`** (from the repo root)

That is the single source of truth — invoke it with `/atempo-design` (or it surfaces
automatically when you're building any Atempo UI, mock, or visual asset). It carries the
selective reading map, the token/component index, the non-negotiables, and a
serve→render→screenshot verification loop.

The rules themselves live in `PRINCIPLES.md` (the invariants + voice) and `README.md` (the
full spec). This file is just a signpost so the pointer survives if the system is copied out
of the repo on its own.
