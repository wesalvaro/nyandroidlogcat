# Nyandroidlogcat

A simple profile-based logcat filterer/highlighter.

## Usage

Configuration profiles can be saved to the file: `~/.nyandroidlogcat.json`
Run with your selected profile: `./nyandroidlogcat profileName`

### Profile Format

```json
{
	"default": { // The profile name (default is used if one is not specified)
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
