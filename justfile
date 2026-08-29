mod api 'bjack-api/api.just'
mod ui 'bjack-ui/ui.just'

default:
  just --list

[group("all")]
setup: api::setup ui::setup

[group("all")]
lint: api::lint ui::lint

[group("all")]
fmt: api::fmt ui::fmt

[group("all")]
fmt_fix: api::fmt_fix ui::fmt_fix

