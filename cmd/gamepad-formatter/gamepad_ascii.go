package main

// Sample gamepad made with https://kenney-assets.itch.io/input-prompts-pixel-16

const inverted_gamepad_unicode = `
      ████████████████████████                                                                              ████████████████████████
    ████████████████████████████                                                                          ████████████████████████████
  ██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████    ████████████████████████████      ████████████████████████████    ██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████
  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████  ████████████████████████████████  ████████████████████████████████  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████
  ████▒▒▒▒▒▒▒▒▓▓▒▒▒▒▓▓▓▓▓▓▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████  ████▒▒▒▒▓▓▓▓▓▓▒▒▓▓▓▓▓▓▒▒▒▒▒▒████
  ████▒▒▒▒▒▒▒▒▓▓▒▒▒▒▒▒▓▓▒▒▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▓▓▒▒▒▒▓▓▓▓▒▒▒▒████  ████▒▒▓▓▓▓▓▓▒▒▓▓▓▓▒▒▒▒▒▒▒▒▒▒████  ████▒▒▒▒▓▓▒▒▓▓▒▒▒▒▓▓▒▒▒▒▒▒▒▒████
  ████░░▒▒▒▒▒▒▓▓▒▒▒▒▒▒▓▓▒▒▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▓▓▒▒▒▒▓▓▓▓▒▒▒▒████  ████▒▒▓▓▒▒▓▓▒▒▓▓▓▓▒▒▒▒▒▒▒▒▒▒████  ████▒▒▒▒▓▓▓▓▒▒▒▒▒▒▓▓▒▒▒▒▒▒░░████
  ████▓▓▒▒▒▒▒▒▓▓▓▓▓▓▒▒▓▓▒▒▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▓▓▒▒▒▒▓▓▒▒▓▓▒▒████  ████▒▒▓▓▓▓▒▒▒▒▓▓▒▒▓▓▒▒▒▒▒▒▒▒████  ████▒▒▒▒▓▓▒▒▓▓▒▒▒▒▓▓▒▒▒▒▒▒▓▓████
  ████▓▓░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████  ████░░▒▒▒▒▒▒▒▒▓▓▓▓▒▒▓▓▓▓▓▓▒▒████  ████▒▒▓▓▒▒▓▓▒▒▓▓▓▓▓▓▒▒▒▒▒▒░░████  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░▓▓████
  ██████▓▓░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░████  ████▓▓░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░▓▓████  ████░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░▓▓██████
    ████▓▓▓▓░░░░░░░░░░░░░░░░░░████  ██████▓▓▓▓░░░░░░░░░░░░░░░░░░████  ████░░░░░░░░░░░░░░░░░░▓▓▓▓██████  ████░░░░░░░░░░░░░░░░░░▓▓▓▓████              ████████████
    ██████▓▓▓▓▓▓░░░░░░░░░░░░▓▓████    ████████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████  ████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████████    ████▓▓░░░░░░░░░░░░▓▓▓▓▓▓██████          ████████████████████
      ██████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████      ████████████████████████████  ████████████████████████████      ████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██████          ████████▒▒▒▒▒▒▒▒████████
        ████████▓▓▓▓▓▓▓▓▓▓▓▓██████          ██████████████████████      ██████████████████████          ██████▓▓▓▓▓▓▓▓▓▓▓▓████████            ████▒▒░░░░▒▒▒▒░░░░▒▒████
          ██████████████████████                                                                          ██████████████████████            ██████▒▒▓▓▓▓░░░░▓▓▓▓▒▒██████
              ████████████████                                                                              ████████████████                ████▒▒▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒▒▒████
                                                                                                                                            ████▒▒▒▒▒▒░░▓▓▓▓░░▒▒▒▒▒▒████
                                                                                                                                            ████▒▒▒▒░░▓▓▓▓▓▓▓▓░░▒▒▒▒████
          ████████████████                                                                                        ████████████              ████▓▓▒▒▓▓▓▓▒▒▒▒▓▓▓▓▒▒▓▓████
        ████████████████████            ████████████████████████          ████████████████████████            ████████████████████          ██████░░░░▒▒▒▒▒▒▒▒░░░░██████
      ██████▒▒▒▒▒▒▒▒▒▒▒▒██████        ████████████████████████████      ████████████████████████████        ████████▒▒▒▒▒▒▒▒████████          ████▓▓▓▓░░░░░░░░▓▓▓▓████
      ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████      ██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████  ██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████      ████▒▒░░░░▒▒▒▒░░░░▒▒████          ████████▓▓▓▓▓▓▓▓████████
      ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████      ████░░▒▒▒▒▒▒▒▒▒▒▓▓▒▒▒▒▒▒▒▒░░████  ████░░▒▒▒▒▒▒▒▒▓▓▒▒▒▒▒▒▒▒▒▒░░████    ██████▒▒▓▓▓▓░░░░▓▓▓▓▒▒██████          ████████████████████
      ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████      ████░░▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒░░████  ████░░▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒░░████    ████▒▒▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒▒▒████              ████████████
      ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████      ████░░▒▒▒▒▒▒▓▓▓▓▓▓▒▒▒▒▒▒▒▒░░████  ████░░▒▒▒▒▒▒▒▒▓▓▓▓▓▓▒▒▒▒▒▒░░████    ████▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒████
      ████░░▒▒▒▒▒▒▒▒▒▒▒▒░░████      ████░░▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒░░████  ████░░▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒░░████    ████▒▒▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▒▒████
      ████▓▓░░░░░░░░░░░░▓▓████      ████░░▒▒▒▒▒▒▒▒▒▒▓▓▒▒▒▒▒▒▒▒░░████  ████░░▒▒▒▒▒▒▒▒▓▓▒▒▒▒▒▒▒▒▒▒░░████    ████▓▓▒▒▒▒▒▒▓▓▓▓▒▒▒▒▒▒▓▓████
      ██████▓▓▓▓▓▓▓▓▓▓▓▓██████      ██████░░░░░░░░░░░░░░░░░░░░██████  ██████░░░░░░░░░░░░░░░░░░░░██████    ██████░░░░▒▒▒▒▒▒▒▒░░░░██████                              ████████████
        ████████████████████          ████████████████████████████      ████████████████████████████        ████▓▓▓▓░░░░░░░░▓▓▓▓████                            ████████████████████
          ████▓▓▓▓▓▓▓▓████              ████████████████████████          ████████████████████████          ████████▓▓▓▓▓▓▓▓████████                          ████████▒▒▒▒▒▒▒▒████████
          ████████████████                                                                                    ████████████████████                            ████▒▒▒▒▒▒░░░░▒▒▒▒▒▒████
            ████████████                                                                                          ████████████                              ██████▒▒▒▒░░▓▓▓▓░░▒▒▒▒██████
                                                                                                                                                            ████▒▒▒▒░░▓▓▓▓▓▓▓▓░░▒▒▒▒████
                                                                                                                                                            ████▒▒▒▒▓▓▓▓░░░░▓▓▓▓▒▒▒▒████
                            ████████████                                                                                                                    ████▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒▒████
                          ████████████████                                    ████████████████                                    ████████████              ████▓▓▒▒▓▓▓▓▒▒▒▒▓▓▓▓▒▒▓▓████
                          ████▒▒▒▒▒▒▒▒████                                  ████████████████████                              ████████████████████          ██████░░░░▒▒▒▒▒▒▒▒░░░░██████
                          ████▒▒▒▒▒▒▒▒████                                ██████▒▒▒▒▒▒▒▒▒▒▒▒██████                          ████████▒▒▒▒▒▒▒▒████████          ████▓▓▓▓░░░░░░░░▓▓▓▓████
                    ██████████▒▒▒▒▒▒▒▒██████████                          ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                          ████▒▒░░░░░░░░░░▒▒▒▒████          ████████▓▓▓▓▓▓▓▓████████
                  ████████████▒▒▒▒▒▒▒▒████████████                        ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                        ██████▒▒▓▓▓▓▓▓▓▓▓▓░░▒▒██████          ████████████████████
                  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                        ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                        ████▒▒▒▒▓▓▓▓░░░░▓▓▓▓▒▒▒▒████              ████████████
                  ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                        ████▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒████                        ████▒▒▒▒▓▓▓▓▓▓▓▓▓▓░░▒▒▒▒████
                  ████░░░░░░░░▒▒▒▒▒▒▒▒░░░░░░░░████                        ████░░▒▒▒▒▒▒▒▒▒▒▒▒░░████                        ████▒▒▒▒▓▓▓▓░░░░▓▓▓▓▒▒▒▒████
                  ████▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▓▓▓▓▓▓▓▓████                        ████▓▓░░░░░░░░░░░░▓▓████                        ████▓▓▒▒▓▓▓▓▓▓▓▓▓▓▒▒▒▒▓▓████
                  ████████████▒▒▒▒▒▒▒▒████████████                        ██████▓▓▓▓▓▓▓▓▓▓▓▓██████                        ██████░░░░▒▒▒▒▒▒▒▒░░░░██████
                    ██████████▒▒▒▒▒▒▒▒██████████                            ████████████████████                            ████▓▓▓▓░░░░░░░░▓▓▓▓████
                          ████░░░░░░░░████                                    ████▓▓▓▓▓▓▓▓████                              ████████▓▓▓▓▓▓▓▓████████
                          ████▓▓▓▓▓▓▓▓████                                    ████████████████                                ████████████████████
                          ████████████████                                      ████████████                                      ████████████
                            ████████████
`

const gamepad_unicode = `
      ░░░░░░░░░░░░░░░░░░░░░░░░                                                                              ░░░░░░░░░░░░░░░░░░░░░░░░
    ░░░░░░░░░░░░░░░░░░░░░░░░░░░░                                                                          ░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░    ░░░░░░░░░░░░░░░░░░░░░░░░░░░░      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░    ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░
  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░
  ░░░░▓▓▓▓▓▓▓▓▒▒▓▓▓▓▒▒▒▒▒▒▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▒▒▒▒▒▒▓▓▒▒▒▒▒▒▓▓▓▓▓▓░░░░
  ░░░░▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▒▒▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▒▒▒▒▓▓▓▓░░░░  ░░░░▓▓▒▒▒▒▒▒▓▓▒▒▒▒▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▒▒▓▓▒▒▓▓▓▓▒▒▓▓▓▓▓▓▓▓░░░░
  ░░░░██▓▓▓▓▓▓▒▒▓▓▓▓▓▓▒▒▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▒▒▒▒▓▓▓▓░░░░  ░░░░▓▓▒▒▓▓▒▒▓▓▒▒▒▒▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▒▒▒▒▓▓▓▓▓▓▒▒▓▓▓▓▓▓██░░░░
  ░░░░▒▒▓▓▓▓▓▓▒▒▒▒▒▒▓▓▒▒▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▒▒▓▓▒▒▓▓░░░░  ░░░░▓▓▒▒▒▒▓▓▓▓▒▒▓▓▒▒▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▒▒▓▓▒▒▓▓▓▓▒▒▓▓▓▓▓▓▒▒░░░░
  ░░░░▒▒██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▓▓▒▒▒▒▒▒▓▓░░░░  ░░░░▓▓▒▒▓▓▒▒▓▓▒▒▒▒▒▒▓▓▓▓▓▓██░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██▒▒░░░░
  ░░░░░░▒▒██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██░░░░  ░░░░▒▒████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓████▒▒░░░░  ░░░░██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██▒▒░░░░░░
    ░░░░▒▒▒▒██████████████████░░░░  ░░░░░░▒▒▒▒██████████████████░░░░  ░░░░██████████████████▒▒▒▒░░░░░░  ░░░░██████████████████▒▒▒▒░░░░              ░░░░░░░░░░░░
    ░░░░░░▒▒▒▒▒▒████████████▒▒░░░░    ░░░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░  ░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░    ░░░░▒▒████████████▒▒▒▒▒▒░░░░░░          ░░░░░░░░░░░░░░░░░░░░
      ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░      ░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░          ░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░
        ░░░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░          ░░░░░░░░░░░░░░░░░░░░░░      ░░░░░░░░░░░░░░░░░░░░░░          ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░░░            ░░░░▓▓████▓▓▓▓████▓▓░░░░
          ░░░░░░░░░░░░░░░░░░░░░░                                                                          ░░░░░░░░░░░░░░░░░░░░░░            ░░░░░░▓▓▒▒▒▒████▒▒▒▒▓▓░░░░░░
              ░░░░░░░░░░░░░░░░                                                                              ░░░░░░░░░░░░░░░░                ░░░░▓▓▓▓▓▓▒▒▒▒▒▒▒▒▓▓▓▓▓▓░░░░
                                                                                                                                            ░░░░▓▓▓▓▓▓██▒▒▒▒██▓▓▓▓▓▓░░░░
                                                                                                                                            ░░░░▓▓▓▓██▒▒▒▒▒▒▒▒██▓▓▓▓░░░░
          ░░░░░░░░░░░░░░░░                                                                                        ░░░░░░░░░░░░              ░░░░▒▒▓▓▒▒▒▒▓▓▓▓▒▒▒▒▓▓▒▒░░░░
        ░░░░░░░░░░░░░░░░░░░░            ░░░░░░░░░░░░░░░░░░░░░░░░          ░░░░░░░░░░░░░░░░░░░░░░░░            ░░░░░░░░░░░░░░░░░░░░          ░░░░░░████▓▓▓▓▓▓▓▓████░░░░░░
      ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░        ░░░░░░░░░░░░░░░░░░░░░░░░░░░░      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░        ░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░          ░░░░▒▒▒▒████████▒▒▒▒░░░░
      ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░      ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░  ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░      ░░░░▓▓████▓▓▓▓████▓▓░░░░          ░░░░░░░░▒▒▒▒▒▒▒▒░░░░░░░░
      ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░      ░░░░██▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▓▓██░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▓▓▓▓██░░░░    ░░░░░░▓▓▒▒▒▒████▒▒▒▒▓▓░░░░░░          ░░░░░░░░░░░░░░░░░░░░
      ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░      ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓██░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓██░░░░    ░░░░▓▓▓▓▓▓▒▒▒▒▒▒▒▒▓▓▓▓▓▓░░░░              ░░░░░░░░░░░░
      ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░      ░░░░██▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓▓▓██░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▒▒▓▓▓▓▓▓██░░░░    ░░░░▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓░░░░
      ░░░░██▓▓▓▓▓▓▓▓▓▓▓▓██░░░░      ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓██░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓██░░░░    ░░░░▓▓▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓░░░░
      ░░░░▒▒████████████▒▒░░░░      ░░░░██▓▓▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▓▓██░░░░  ░░░░██▓▓▓▓▓▓▓▓▒▒▓▓▓▓▓▓▓▓▓▓██░░░░    ░░░░▒▒▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▒▒░░░░
      ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░      ░░░░░░████████████████████░░░░░░  ░░░░░░████████████████████░░░░░░    ░░░░░░████▓▓▓▓▓▓▓▓████░░░░░░                              ░░░░░░░░░░░░
        ░░░░░░░░░░░░░░░░░░░░          ░░░░░░░░░░░░░░░░░░░░░░░░░░░░      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░        ░░░░▒▒▒▒████████▒▒▒▒░░░░                            ░░░░░░░░░░░░░░░░░░░░
          ░░░░▒▒▒▒▒▒▒▒░░░░              ░░░░░░░░░░░░░░░░░░░░░░░░          ░░░░░░░░░░░░░░░░░░░░░░░░          ░░░░░░░░▒▒▒▒▒▒▒▒░░░░░░░░                          ░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░
          ░░░░░░░░░░░░░░░░                                                                                    ░░░░░░░░░░░░░░░░░░░░                            ░░░░▓▓▓▓▓▓████▓▓▓▓▓▓░░░░
            ░░░░░░░░░░░░                                                                                          ░░░░░░░░░░░░                              ░░░░░░▓▓▓▓██▒▒▒▒██▓▓▓▓░░░░░░
                                                                                                                                                            ░░░░▓▓▓▓██▒▒▒▒▒▒▒▒██▓▓▓▓░░░░
                                                                                                                                                            ░░░░▓▓▓▓▒▒▒▒████▒▒▒▒▓▓▓▓░░░░
                            ░░░░░░░░░░░░                                                                                                                    ░░░░▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▓▓▓▓░░░░
                          ░░░░░░░░░░░░░░░░                                    ░░░░░░░░░░░░░░░░                                    ░░░░░░░░░░░░              ░░░░▒▒▓▓▒▒▒▒▓▓▓▓▒▒▒▒▓▓▒▒░░░░
                          ░░░░▓▓▓▓▓▓▓▓░░░░                                  ░░░░░░░░░░░░░░░░░░░░                              ░░░░░░░░░░░░░░░░░░░░          ░░░░░░████▓▓▓▓▓▓▓▓████░░░░░░
                          ░░░░▓▓▓▓▓▓▓▓░░░░                                ░░░░░░▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░                          ░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░          ░░░░▒▒▒▒████████▒▒▒▒░░░░
                    ░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░                          ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                          ░░░░▓▓██████████▓▓▓▓░░░░          ░░░░░░░░▒▒▒▒▒▒▒▒░░░░░░░░
                  ░░░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░░░                        ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                        ░░░░░░▓▓▒▒▒▒▒▒▒▒▒▒██▓▓░░░░░░          ░░░░░░░░░░░░░░░░░░░░
                  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                        ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                        ░░░░▓▓▓▓▒▒▒▒████▒▒▒▒▓▓▓▓░░░░              ░░░░░░░░░░░░
                  ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                        ░░░░▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░                        ░░░░▓▓▓▓▒▒▒▒▒▒▒▒▒▒██▓▓▓▓░░░░
                  ░░░░████████▓▓▓▓▓▓▓▓████████░░░░                        ░░░░██▓▓▓▓▓▓▓▓▓▓▓▓██░░░░                        ░░░░▓▓▓▓▒▒▒▒████▒▒▒▒▓▓▓▓░░░░
                  ░░░░▒▒▒▒▒▒▒▒▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒░░░░                        ░░░░▒▒████████████▒▒░░░░                        ░░░░▒▒▓▓▒▒▒▒▒▒▒▒▒▒▓▓▓▓▒▒░░░░
                  ░░░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░░░                        ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░                        ░░░░░░████▓▓▓▓▓▓▓▓████░░░░░░
                    ░░░░░░░░░░▓▓▓▓▓▓▓▓░░░░░░░░░░                            ░░░░░░░░░░░░░░░░░░░░                            ░░░░▒▒▒▒████████▒▒▒▒░░░░
                          ░░░░████████░░░░                                    ░░░░▒▒▒▒▒▒▒▒░░░░                              ░░░░░░░░▒▒▒▒▒▒▒▒░░░░░░░░
                          ░░░░▒▒▒▒▒▒▒▒░░░░                                    ░░░░░░░░░░░░░░░░                                ░░░░░░░░░░░░░░░░░░░░
                          ░░░░░░░░░░░░░░░░                                      ░░░░░░░░░░░░                                      ░░░░░░░░░░░░
                            ░░░░░░░░░░░░
`