FROM ubuntu:20.04

RUN apt-get update && \
	DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get install -y \
		build-essential \
		libsdl2-dev \
		mingw-w64 \
		curl \
		unzip && \
	bash -c 'echo "SDL_URL=https://www.libsdl.org" > /tmp/download-sdl.sh' && \
	bash -c "echo -n 'RELEASE_ENDPOINT=' >> /tmp/download-sdl.sh" && \
		bash -c "echo -n '\`' >> /tmp/download-sdl.sh" && \
		bash -c "echo -n 'curl \"\${SDL_URL}/download-2.0.php\" | ' >> /tmp/download-sdl.sh" && \
		bash -c "echo -n \"grep 'release.*mingw.*tar.gz' | \" >> /tmp/download-sdl.sh" && \
		bash -c "echo -n \"awk -v FS=\\\\\\\" '{print \\\$2}'\" >> /tmp/download-sdl.sh" && \
		bash -c "echo '\`' >> /tmp/download-sdl.sh" && \
	bash -c "echo 'curl --output /tmp/sdl-mingw.tar.gz \"\${SDL_URL}/\${RELEASE_ENDPOINT}\"' >> /tmp/download-sdl.sh" && \
	bash /tmp/download-sdl.sh && \
	tar -zxf /tmp/sdl-mingw.tar.gz -C /tmp/ && \
	tar -cC /tmp/SDL2-2.0.22/i686-w64-mingw32/ `ls /tmp/SDL2-2.0.22/i686-w64-mingw32/` | tar -xC /usr/i686-w64-mingw32/ && \
	tar -cC /tmp/SDL2-2.0.22/x86_64-w64-mingw32/ `ls /tmp/SDL2-2.0.22/x86_64-w64-mingw32/` | tar -xC /usr/x86_64-w64-mingw32/ && \
	rm -rf /tmp/SDL2-2.0.22 && \
	rm -f /tmp/download-sdl.sh
