# Bismuth

A logic circuit simulator built in Go. Place wires, gates, and components on a grid, toggle switches, and watch electricity flow in real time.

---

## How it works

The world is a 3D grid of tiles: X columns, Y rows, and Z layers stacked on top of each other. Every tile on the grid holds one component, a power state, and a direction it faces.

Each frame, the engine runs a three-phase power update:

1. **Read** - scans every tile to find active power sources (switches that are on, buttons with time left, NOT gates whose input is off)
2. **Wipe** - resets the power state of every wire, light, diode, and via to unpowered
3. **Spread** - starts from every active source and floods outward through connected wires, respecting diode directions and jumping between layers through vias

This means the circuit is fully re-evaluated every frame. There is no state carried over from the previous tick except for switches, buttons, and the tiles whose power is determined by their own logic (like NOT gates).

---



## Controls

Press `H` inside the app to open the help panel. All controls are listed there.

---

## Layers

The grid has 3 layers stacked on top of each other. You can only see and edit one layer at a time. To connect power between layers, place a **Via** on the same X/Y position on both layers - power will jump across automatically.

This lets you route wires over each other without them interfering, which is useful for building more complex circuits.

---

## Getting Started

```bash
go run main.go
```

Requires Go and [Ebitengine](https://ebitengine.org/). The window opens at 800x800 (40 tiles wide, 40 tiles tall, 20px per tile). Three layers are available by default.

---

## Project Structure

```
bismuth/
  main.go                  - window, input handling, rendering
  internal/engine/
    grid.go                - grid, power simulation (read/wipe/spread)
    tile.go                - Tile struct, ComponentType, Direction
```
