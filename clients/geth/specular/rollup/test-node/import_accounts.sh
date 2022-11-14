#!/bin/bash
../../build/bin/geth --password ./password.txt --datadir ./data account import ./keys/sequencer.prv
../../build/bin/geth --password ./password.txt --datadir ./data account import ./keys/validator.prv
../../build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/sequencer.prv
../../build/bin/geth --password ./password.txt --datadir ./data_validator account import ./keys/validator.prv
