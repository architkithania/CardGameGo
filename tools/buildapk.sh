rm -rf android/src/main/assets

cp -rf ../../assets android/src/main/assets

./gradlew assembleDebug
