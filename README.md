# A simple yet effective weather CLI 

`weather` is an simple weather tool in the terminal.
Add your Met Office API key to return real time wheather data to your exact location!

### About
The CLI is currently configured to always return the weather predicted over the rest of the current day in three hour intervals. The last interval available from the Met Office API is 21:00, so if the CLI is used after this time, the most recent entry will be all that is displayed. 

<p align="center">
  <img src="./resources/media/swappy-20231202-124415.png?raw=true" alt="Example1" />
</p>

<p align="center">
  <img src="./resources/media/swappy-20231202-124719.png?raw=true" alt="Example2" />
</p>


This will change inline with your terminal's theme. If any of the icons do not parse, you may need to install a font that is compatible with color/emojis such as [NerdFonts](https://docs.rockylinux.org/books/nvchad/nerd_fonts/)

## Installation

```
go install github.com/savedra1/weather/weather@latest
```
This will download the executable to `${HOME}/go/bin`. Move the `weather` binary file from here to your local `/bin` / `/usr/bin` / `/.local/bin` or set up a custom alias so that the `weather` command can be ran from any dir. 

### MacOS
On MacOS, the entire install can be completed with the following command: 
```
sudo export GOBIN=/usr/local/bin && sudo go install github.com/MichaelSavedra1/weather/weather@latest 
```
### Windows
The main installation command will download a .exe file that can be set to a CMD alias/ran manually. 

## Usage after installing
1. Create a free [Met Office Data Point account](https://register.metoffice.gov.uk/MyAccountClient/account/view) (easy sign-up)
2. Copy your application key from the [My Account](https://register.metoffice.gov.uk/MyAccountClient/account/view) page 
3. Run the following command in your terminal and you will be prompte to supply your application key
```
weather
```
4. Use `weather --help` to see a list of all available commands
5. Set your default location by using the `weather set-default` arg. To see a full list of locations, check the resources/met-api directory
6. Use `weather --extended` or `weather {arg} --extended` to see the next 5 days worth of met office weather data

### License

MIT