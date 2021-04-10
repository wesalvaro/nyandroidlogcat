# Nyandroidlogcat

A simple profile-based logcat filterer/highlighter.

## Usage

There are two different _breeds_ of logcat for you:

  - **Stream-based** \
    Outputs the logs based on your configured profiles to standard out. \
    There is no interaction with the stream, scrolling is from your terminal. \
    This version is fairly robust in terms of output.

    ```shell
    cd streamcat && go build # Build
    ./streamcat # Run
    ```

    Configuration profiles can be saved to the file: `~/.nyandroidlogcat.json` \
    Run with a non-default profile: `./nyandroidlogcat profileName`

  - **Curses-based** \
    Uses TermUI to create interactive UI in your terminal. \
    Allows for interactive filtering _as the logs are coming in_. \
    When the log is slow (or filtered to be slow) this version's interactive
    filtering is quite nice and convenient. However, when the log is flying by,
    this version may experience glitches and scrolling is impossible.

    ```shell
    cd cursedcat && go build # Build
    ./cursedcat # Run
    ```

    - Type to filter the log message
    - Use left/right to set the lower bound filter for the log level
    - Use up/down, home/end, ctrl+j/k to move around
    - Press escape or ctrl+c to quit

### Profile Format

```js
{
  "default": { // The profile name ("default" is used if one is not specified)
    "time-format": "15:04:05", // Go time layout (Jan 2 15:04:05 2006 MST)
    "tag": { // Options for the logcat tag
      "show": true, // Show tags in the output
      "filter": [], // Show only tags matched by these regular expressions
      "ignore": [] // Ignore tags matched by these regular expressions
    },
    "message": { // Options for the logcat message
      "filter": [], // Show only messages matches by these regular expressions
      "highlight": [], // Highlight these regular expressions in the output
      "highlight-color": "yellow" // Background color for the highlight
    },
    "level": { // Options for the logcat level
      "show": true, // Show levels in the output
      "bound": "Info", // The lower-bound of the level in the output
      "first": true, // Show the level first on the output line
      "color": true, // Use a different color for each level
      "emoji": true, // Use an emoji instead of ASCII for each level
      "long": false // Print the full level word (e.g. "Warning")
    }
  }
}
```
