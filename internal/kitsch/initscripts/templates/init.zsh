# Adapted from https://github.com/starship/starship/blob/master/src/init/starship.zsh
# Copyright (c) 2019-2021, Starship Contributors

# ZSH has a quirk where `preexec` is only run if a command is actually run (i.e
# pressing ENTER at an empty command line will not cause preexec to fire). This
# can cause timing issues, as a user who presses "ENTER" without running a command
# will see the time to the start of the last command, which may be very large.

# To fix this, we create KITSCH_START_TIME upon preexec() firing, and destroy it
# after drawing the prompt. This ensures that the timing for one command is only
# ever drawn once (for the prompt immediately after it is run).

zmodload zsh/parameter  # Needed to access jobstates variable for KITSCH_JOBS_COUNT

# If you try to install kitsch prompt overtop older versions of starship, they
# clash.  Remove any starship config if it's there.
precmd_functions[$precmd_functions[(i)starship_precmd]]=()
preexec_functions[$preexec_functions[(i)preexec_functions]]=()
starship_zle-keymap-select() {}

# Defines a function `__kitschprompt_get_time` that sets the time since epoch in millis in KITSCH_CAPTURED_TIME.
if [[ $ZSH_VERSION == ([1-4]*) ]]; then
    # ZSH <= 5; Does not have a built-in variable so we will rely on Starship's inbuilt time function.
    __kitschprompt_get_time() {
        KITSCH_CAPTURED_TIME=$("{{ .kitschCommand }}" time)
    }
else
    zmodload zsh/datetime
    zmodload zsh/mathfunc
    __kitschprompt_get_time() {
        (( KITSCH_CAPTURED_TIME = int(rint(EPOCHREALTIME * 1000)) ))
    }
fi

# Will be run before every prompt draw
kitsch_precmd() {
    # Save the status, because commands in this pipeline will change $?
    KITSCH_CMD_STATUS=$?

    # Compute cmd_duration, if we have a time to consume, otherwise clear the
    # previous duration
    if (( ${+KITSCH_START_TIME} )); then
        __kitschprompt_get_time && (( KITSCH_DURATION = KITSCH_CAPTURED_TIME - KITSCH_START_TIME ))
        unset KITSCH_START_TIME
    else
        unset KITSCH_DURATION
    fi

    # Use length of jobstates array as number of jobs. Expansion fails inside
    # quotes so we set it here and then use the value later on.
    KITSCH_JOBS_COUNT=${#jobstates}
}
kitsch_preexec() {
    __kitschprompt_get_time && KITSCH_START_TIME=$KITSCH_CAPTURED_TIME
}

# If precmd/preexec arrays are not already set, set them. If we don't do this,
# the code to detect whether kitsch_precmd is already in precmd_functions will
# fail because the array doesn't exist (and same for kitsch_preexec)
(( ! ${+precmd_functions} )) && precmd_functions=()
(( ! ${+preexec_functions} )) && preexec_functions=()

# If starship precmd/preexec functions are already hooked, don't double-hook them
# to avoid unnecessary performance degradation in nested shells
if [[ -z ${precmd_functions[(re)kitsch_precmd]} ]]; then
    precmd_functions+=(kitsch_precmd)
fi
if [[ -z ${preexec_function[(re)kitsch_preexec]} ]]; then
    preexec_functions+=(kitsch_preexec)
fi

# Set up a function to redraw the prompt if the user switches vi modes
kitsch_zle-keymap-select() {
    zle reset-prompt
}

## Check for existing keymap-select widget.
# zle-keymap-select is a special widget so it'll be "user:fnName" or nothing. Let's get fnName only.
__kitsch_preserved_zle_keymap_select=${widgets[zle-keymap-select]#user:}
if [[ -z $__kitsch_preserved_zle_keymap_select ]]; then
    zle -N zle-keymap-select kitsch_zle-keymap-select;
else
    # Define a wrapper fn to call the original widget fn and then Starship's.
    kitsch_zle-keymap-select-wrapped() {
        $__kitsch_preserved_zle_keymap_select "$@";
        kitsch_zle-keymap-select "$@";
    }
    zle -N zle-keymap-select kitsch_zle-keymap-select-wrapped;
fi

__kitschprompt_get_time && KITSCH_START_TIME=$KITSCH_CAPTURED_TIME

# Set up the session key that will be used to store logs
KITSCH_SESSION_KEY="$RANDOM$RANDOM$RANDOM$RANDOM$RANDOM"; # Random generates a number b/w 0 - 32767
KITSCH_SESSION_KEY="${KITSCH_SESSION_KEY}0000000000000000" # Pad it to 16+ chars.
export KITSCH_SESSION_KEY=${KITSCH_SESSION_KEY:0:16}; # Trim to 16-digits if excess.

VIRTUAL_ENV_DISABLE_PROMPT=1

setopt promptsubst
PROMPT='$("{{ .kitschCommand }}" prompt {{with .configFile}}--config {{.}} {{end}}--shell zsh --terminal-width="$COLUMNS" --keymap="$KEYMAP" --status="$KITSCH_CMD_STATUS" --cmd-duration="$KITSCH_DURATION" --jobs="$KITSCH_JOBS_COUNT")'
