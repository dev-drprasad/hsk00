- Console is able to extract `.asd` files with 6.5 MB
- seems like rom size should be less than 1 MB for to play
- seen problem with multi-cart roms
- random names for `.asd` and `.bin` files works as long as `hsk00.asd` file is updated
- ncs files are just like asd files. `descramble` command can extract game from `.ncs` file
- Menu.logXX are also images. `Menu.log07` to `Menu.log25` contains mushroom selector in menu
- `.logXX`, `.binXX` are interchangeble
- XOR of `c8` and value at offset: `04` (.bin, logXX) will give image starting first two bytes
- In "(1400 + 520 rus) model", hsk00.asd contains animated mushroom slices, but in "Y2 model", hsk00.asd contains game list
- `Menu.ocv` (both Y2 and Y2 plus) contains background music
-
- `classics_menu_ok.drm` (Not sure from which console) and `Menu.ocv` (from Y2) have same music. But their content is different though
- Files that starts with `SP_ToneMaker`/`SPF2ALP` contains music/sound
- I assumed `Menu.ocv` contains PCM raw data and played with audacity (8bit unsinged, Mono, saample rate: 22050), I hear music, but there is too much noise. Either my settings are wrong or its not PCM
- G+ gadget audio file batch converter can be used to make "data frog Y2" supported music.
  - [Download](https://www.generalplus.com/1LVlangLNxxSVyySNservice_n_support_d)
  - Alogirthm: H/W PCM 16-Bit (DRM)
  - I tried with sample rates: 22050, 11025 wav input files
  - Output drm file more than 3MB didn't work for me

### REFERENCES

- https://git.redump.net/mame/patch/?id=9bf6c963596e16381fd793e729420fdc14c7e326
- https://ia802802.us.archive.org/20/items/gpac800/GPL162004A-507A_162005A-707AV10_code_reference-20147131205102.pdf
- https://forums.bannister.org/ubbthreads.php?ubb=showflat&Number=112652#Post112652
- https://www.generalplus.com/1LVlangLNxxSVyySNservice_n_support_d
- https://github.com/alito/mamele/commit/d4ce622200de7a3b2088524e8798324c890b781b
- https://github.com/search?p=2&q=adpcm36+decoding&type=Code
- https://github.com/Jeija/bluefluff
- https://sites.google.com/site/devusb/waveformat
- http://webcache.googleusercontent.com/search?q=cache:QPxE2Hil8_wJ:bootleg.games/BGC_Forum/index.php%3Faction%3Dprofile%3Barea%3Dshowposts%3Bsa%3Dtopics%3Bu%3D38
- online-convert.com
- https://hackmii.com/2010/04/sunplus-the-biggest-chip-company-youve-never-heard-of/
- https://www.dl-sounds.com/royalty-free

### SD Card Backups

- https://yadi.sk/d/rTNU4YlXHUIRyg
- Y2 HD: https://4pda.ru/forum/dl/post/21973244/Y2-HD+-+games-EN.zip
- Y2 HD (alternate link): https://disk.yandex.ru/d/JMBiacwqIiktOw?w=1
- Y3: https://yadi.sk/d/POgSN-oMOxI-hw
- Y2S 750+ https://drive.google.com/uc?id=1qKgWlQhW7WYfkd_vzOH1op5E5dfQwMEw&export=download
