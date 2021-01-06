# hsk00

Add/Replace games to your "Data Frog Y2 (568 in 1)" console. This is CLI tool. Download tool from [releases](https://github.com/dev-drprasad/hsk00/releases)

🚧 **This is work in progress** 🚧

---

<br />

## API

⚠️ BEFORE YOU RUN ANY COMMAND, BACKUP YOUR SD CARD ⚠️
<br /><br />

### `add`

Adds game(s) to given category (Sports, Adventure etc..)

`--category` number starts from 0, left to right in menu. Example: "Sports Games" category number is `4`.

`--root` is root directory of game folder where `Menu.ocv` exists (can be sd card path or custom directory where files present)

**Example:**

```shell
hsk00-darwin-amd64 add  in/Famicom_Wars.nes  in/Heavy_Barrel.nes  in/Fantasy_Zone.nes  in/Final_Fantasy_II.nes --category 4 --root ~/Datafrog
```

⚠️ This will change menu text slightly ⚠️. If you can't lanuch games or change menu page, restore files with your backup.

---

<br />

### `replace`

Replace will replace existing game with custom game.

Not implemented yet

---

<br />

### `descramble`

Converts `*.asd` files to usable `.zip` files. Output filename will be `<inputfilename>.zip` and will be generated in same directory where input file present.

Example:

```
hsk00-darwin-amd64 descramble  ~/Datafrog2/Game04/hsk06.asd
```

<br />

## Supported Consoles

- DATA FROG Y2 HD (568 in 1)
- Extreme Mini HD Game Box
- Probably works with SD card that looks like below

<img src="./sd-layout.png"  alt="data-frog-sd-card-files" width="500"  />

Let me know if it works with other consoles. It helps other people

## TODO

- [ ] May be GUI ?
- [ ] replace game
- [ ] delete game

## Need Help

I am not able to understand what are `GameXX.bin`, `Menu.logXX` files. They have background images, and menu selection images hidden them. But I am not able to determine offset of these images. If you know anything about these files, please let me know.

## References

- http://bootleg.games/BGC_Forum/index.php?PHPSESSID=bvomlllrtphq11187kpvontr72&topic=1775.msg17586#msg17586
- https://golangcode.com/create-zip-files-in-go/
- https://gist.github.com/madevelopers/40b269730df687cdcb8b
- https://stackoverflow.com/questions/28513486/how-add-a-file-to-an-existing-zip-file-using-golang
- http://blog.ralch.com/tutorial/golang-working-with-zip/
- https://stackoverflow.com/a/42454716/6748719
- https://exifinfo.org/
