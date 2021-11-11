# templates

Most templates in this folder are stolen from Starship's initialization templates. To adapt one:

- Replace `::STARSHIP::` with `{{ .kitschCommand }}`.
- Replace all instances of "starship" with "kitsch".
- Remvoe the STARSHIP_SHELL environment variable - instead pass this value as `--shell` to the prompt command.
