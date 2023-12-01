# A simple yet effective weather CLI 

`weather` is an simple weather tool in the terminal.
Add your Met Office API key to return real time wheather data to your exact location!

### About
The CLI is currently configured to always return the weather predicted over the rest of the current day in three hour intervals. The last interval available from the Met Office API is 21:00, so if the CLI is used after this time, the most recent entry will be all that is displayed. 

<p align="center">
  <img src="./resources/media/example.png?raw=true" alt="Example" />
</p>

<p align="center">
  <img src="./resources/media/swappy-20231201-194203.png?raw=true" alt="Example" />
</p>


This will change inline with your terminal's theme. If any of the icons do not parse, you may need to install a font that is compatible with color/emojis such as [NerdFonts](https://docs.rockylinux.org/books/nvchad/nerd_fonts/)

*`weather` is currently only compatible with Linux and MacOS machines.*

### Installation

```
go install github.com/MichaelSavedra1/weather@latest
```

### Usage after installing
1. Create a free [Met Office Data Point account](https://register.metoffice.gov.uk/MyAccountClient/account/view) (easy sign-up)
2. Copy your application key from the [My Account](https://register.metoffice.gov.uk/MyAccountClient/account/view) page 
3. Run the following command in your terminal and you will be prompte to supply your application key
```
weather
```
4. Use `weather --help` to see a list of all available commands
5. Set your default location by using the `weather set-default` arg. To see a full list of locations, check the resources/met-api directory

### License

MIT