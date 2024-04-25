# [WIP] SwitchBud

A pack switcher within Discord, built on top of [Crafty Controller](https://craftycontrol.com/).

## Why?

CC (Crafty Controller) is "a free and open-source Minecraft launcher and manager that allows
users to start and administer Minecraft servers from a user-friendly interface." This is neat,
but only gets you so far. This is where SwitchBud comes in. If you're hosting multiple servers,
but don't have the resources to host them all at once, SwitchBud is perfect for you. SwitchBud
acts as a Discord bot that allows you to switch packs on the fly. This is especially useful if
you want others to switch between packs, but don't want to give them access to the admin panel.

## Installation

```sh
go install github.com/tizu69/switchbud@latest
switchbud init
vi switchbud.yml
switchbud run
```

---

these docs are WIP! here's some misc stuff:

## Config / Resource Slots

A resource slot is an arbitrary unit. You can use any value there, although I usually use the
average RAM usage for each pack (but anything goes). SwitchBud will, on pack switch attempt,
check how many slots are available (all servers share the same slot pool). If there are less
slots available than required, the pack won't be able to boot (at least that's what SwitchBud
thinks) - so it won't even try. If you only want a single server running at a time, you can
change the slot count to 1 and change each server's slot usage to 1. SwitchBud is smart enough
to shut down servers to free up just enough slots to run the server you want, but it'll refuse
to do so if players are online.
