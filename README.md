# sensor-stream

## Tech/Framework used

### sensor-stream
<b>Scaffolding</b>
- [GO MOD INIT] (https://blog.golang.org/using-go-modules)

<b>Built with</b>
- [Go 1.6] (https://go.dev/blog/go1.6)
- [Testify - Thou Shalt Write Tests] (https://github.com/stretchr/testify)

## Packaging and running

Commands below must be executed from the project root

### Build and Run

`make run` or `docker-compose up --build -d`

### Stop

`make stop` or `docker-compose down`

## Assumptions

<b>When getting available pilots</b>:

- the distribution is based in the amount of flights a pilot has
  in the weekday of the departure request
- If two pilots have the same amount on flights in a given day,
  the decision is by alphabetical order