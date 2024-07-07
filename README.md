# DRAWJ2D-GOLANG PORT

[![rm1](https://img.shields.io/badge/rM1-supported-green)](https://remarkable.com/store/remarkable)
[![rm2](https://img.shields.io/badge/rM2-supported-green)](https://remarkable.com/store/remarkable-2)
[![Discord](https://img.shields.io/discord/385916768696139794.svg?label=reMarkable&logo=discord&logoColor=ffffff&color=7389D8&labelColor=6A7EC2)](https://discord.gg/ATqQGfu)
[![rM Hacks Discord](https://img.shields.io/discord/1153374327123759104.svg?label=rM%20Hacks&logo=discord&logoColor=ffffff&color=ffb759&labelColor=d99c4c)](https://discord.gg/bgVXW2bchN)


This port is meant to be run on the remarkable tablet a Golang + C implementation of remarkable drawj2d conversion code

## Requirements:
- rm-hacks (enables screenshot feature): https://github.com/mb1986/rm-hacks

- webinterface-onboot: https://github.com/rM-self-serve/webinterface-onboot 


⚠️ Please be sure to have rm-hacks and webinterface-onboot Installed⚠️


1. **Transfrer the client installer tar to the remarkable:**

   then:
   - 
   ```bash
   $remarkable: ~/ tar -xvf drawj2d-rm.tar
   $remarkable: ~/ cd drawj2d-rm
   $remarkable: ~/ ./install.sh

   ```
---

Now you should be able to convert your screenshots to rmlines in 3 sec

## Benchmark:

<img src="remarkablepage/bench/cpu-new-bench-CPROCESSING.prof.svg" alt="Benchmark" width="800" height="800">

