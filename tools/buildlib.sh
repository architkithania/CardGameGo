source tools/envs

CC=arm-linux-androideabi-gcc \
CGO_CFLAGS="-D__ANDROID_API__=16 -I${ANDROID_NDK_HOME}/sysroot/usr/include -I${ANDROID_NDK_HOME}/sysroot/usr/include/arm-linux-androideabi --sysroot=${ANDROID_SYSROOT}" \
CGO_LDFLAGS="-L${ANDROID_NDK_HOME}/sysroot/usr/lib -L${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/ --sysroot=${ANDROID_SYSROOT}" \
CGO_ENABLED=1 GOOS=android GOARCH=arm \
go build -tags static -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" -o=creating_apk/android/android/libs/armeabi-v7a/libexample.so /home/architkithania/Documents/Coding/Experimentation/CardGameGo/src/