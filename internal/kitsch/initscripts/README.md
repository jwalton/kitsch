# templates

Most templates in ./tempaltes are stolen from Starship's initialization templates. To adapt a starship one:

- Replace `::STARSHIP::` with `{{ .kitschCommand }}`.
- Remove the STARSHIP_SHELL environment variable - instead pass this value as `--shell` to the prompt command.
- Replace all instances of "starship" with "kitsch".
