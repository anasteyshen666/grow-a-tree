# Game Design Document: "Grow a Tree"

**Genre:** Survival / Real-time strategy / Tower Defense on a grid.
**Platform:** PC (Windows/Linux/Mac).
**Engine:** No engine. Plain **Go (Golang)** + the **Ebiten** library.
**Visual style:** Minimalist, dark fantasy, bioluminescence (glowing neon roots on a dark background).
**Vibe:** Meditative yet tense. A phase of calm planning suddenly gives way to panic during bug hordes.

---

## 1. Elevator Pitch
"Grow a Tree" is an endless strategy game where you play the mind of an ancient forest. You start with a single Core (a 2×2 square) and must expand a network of roots underground, mining water and energy. Survival demands balance: water sources run dry, winter frosts set in, and hordes of pest bugs of 5 different levels crawl in from every side. Use symbiosis with mushrooms, cut off rotting roots as bait, and grow companion plants to survive as many waves as possible and beat your record.

---

## 2. Resources (Game Economy)
Everything is built on the balance of three resources:

* **Energy (green bar):** Spent on growing new roots, planting plants, and activating defense. Regenerates over time, but **only if there is water**.
* **Water (blue bar):** Spent on regenerating Energy. Mined from sources (blue puddles on the map) that run dry over time. The maximum water reserve increases if you connect the roots of different trees (Cores) into a single network.
* **Seeds (yellow bar):** Accumulate over time. Spent on planting Companion Plants ("pets") and creating new Cores (trees).

---

## 3. Core Mechanics
* **Grid and Growth (LMB):** The field is split into cells. LMB spends energy and grows a root into an empty cell if it touches an existing root or the Core.
* **Cutting and Rot (RMB):** The player can remove their own root at any time, refunding 50% of the energy spent. **Rot** (a brown cell) is left where the root was.
* **Network logistics:** If a bug gnaws through a root in the middle, all roots beyond that point "fall off" (stop providing water and energy) until the gap is patched.
* **Maturation cycle:** Pressing the "Mature" button spends accumulated Seeds. A new tree (Core) appears, to which you must run roots to raise the water limit.

---

## 4. Defense and Symbiosis
* **Mushroom Symbiosis:** Mushroom spores appear randomly on the map. If your root touches a spore, it becomes infected and turns "Mushroom". Bugs crawl over such roots twice as slowly. *Downside:* Mushroom roots cannot merge the water reserves of different trees.
* **Rot as Bait:** Bugs prefer fresh Rot. If there is Rot on the map, bugs reroute and head straight for it. After eating rot, a bug is poisoned (slows down and loses health), then disappears. This lets you lure huge level-5 bugs away from the Core.

---

## 5. Companion Plants ("pets")
Plants are planted for Seeds next to a Core. They do not conduct water, but radiate an **aura** within a 2–3 cell radius:
1. **Battery Flower (yellow):** +20% energy regen rate to all roots in radius.
2. **Water Moss (blue):** Reduces water cost of energy regen by 10% in radius.
3. **Winter Thornbush (white):** Radiates heat. Within its range, water sources do not freeze in winter.

---

## 6. Seasonality (Meta-game)
The game runs through seasons. Season changes are tied to wave numbers (e.g. a season changes every 5 waves).
* **Summer:** Standard rules. The ideal time to expand and accumulate Seeds.
* **Autumn:** Water sources start drying up faster. Active mushroom-spore spawning kicks in.
* **Winter:** Water sources freeze (stop providing water). The player must frantically plant Thornbushes around the remaining puddles or live off reserves stored in tanks. Bugs become more armored.
* **Spring:** Ice melts, rich new water sources appear, but the number of spawning bugs doubles.

---

## 7. Enemies (Endless Waves)
Bugs come in waves. There are 15 seconds of prep between waves.
* **Bug AI:** Levels 1 through 5.
  * *Lvl 1–2:* Fast but weak. Easily killed by roots or rot.
  * *Lvl 3–4:* Medium. Gnaw roots faster, able to cut off large chunks of the network.
  * *Lvl 5 (Boss):* Huge, slow, tons of HP. Requires enormous energy to kill or a clever trap of mushrooms and rot.
* **Killing a bug:** A root can attack a bug that steps onto it, but this quickly drains Energy. The bigger the bug, the more energy it takes to destroy.
* **Defeat:** A bug reaches the Core and eats it. If there is only one Core — Game Over.

---

## 8. Records and Replayability
* The game is endless. The main goal is to survive as many waves as possible.
* A `save.json` file is created in the game folder, recording the maximum wave reached (High Score). On the start screen this record is always shown, pushing the player to beat it.
* *Rogue-lite elements:* Because of the random placement of water sources, mushroom spores, and the directions of the first bug waves, no two runs are alike.

---

## 9. Interface (UI) and Controls
* **Top-left:** Energy, Water, and Seeds bars with numbers.
* **Top-right:** Current Wave (e.g. "Wave 12"), Season (leaf/snowflake icon), and Record.
* **Bottom-center:** The "Mature (Seed)" button and the timer until the next bug attack.
* **Controls:** The whole game is played with one mouse (LMB — grow, RMB — cut/rot) and a couple of hotkeys (Space — pause, R — restart after death).
