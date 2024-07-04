docker run --rm \
    --platform linux/amd64 \
    -e HTTPS_PROXY=http://127.0.0.1:7890 \
    --network host \
    -v $(pwd):/workspace \
    -v $(cd ~/.cache/docker_omega_build_cache && pwd):/root/go \
    -v $(cd ~/.cache/docker_omega_build_cache/build_cache && pwd):/root/.cache/go-build \
    neomega:builder_amd64 \
    make -C /workspace build