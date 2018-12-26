for d in 10ns 2ms 5ms; do
    export EMIT_INTERVAL=$d

    for i in seq 5; do
        echo "$(timeout 10s ./continuous_log_emitter | wc -l) logs for interval $d" &
    done
done
