# Create a directory for mono files if it doesn't exist
mkdir -p mono

# Convert all .flac files to mono and store them in the mono/ directory
for file in *.flac; do
    ffmpeg -i "$file" -ac 1 "mono/$file"
done

echo "Finished converting stereo to mono. Files are stored in 'mono' directory."

# Change directory to mono
cd mono

# Create a directory for processed files if it doesn't exist
mkdir -p processed

# Normalize and compress each mono file, then store it in the processed/ directory
for file in *.flac; do
    ffmpeg -i "$file" \
    -filter_complex "[0:a]loudnorm=I=-16:TP=-1.5:LRA=11[a];[a]acompressor=threshold=-24dB:ratio=4.5:attack=3:release=100[a]" \
    -map "[a]" -c:a flac "processed/$file"
done

echo "Finished processing files. Noise-gated, normalized, and compressed files are stored in 'mono/processed' directory."

# Concatenate all processed files into one
cd processed

input_args=""
file_count=0

for file in *.flac; do
  ((file_count++))

  input_args+="-i \"$file\" "
done

eval "ffmpeg $input_args -filter_complex \"amix=inputs=$file_count\" all.flac"

ffmpeg -i all.flac -filter_complex "silenceremove=stop_periods=-1:stop_duration=1:stop_threshold=-50dB" removed_silent.flac

ffmpeg -i removed_silent.flac -filter_complex "acompressor" all_voices_final.flac

ffmpeg -i all_voices_final.flac -ab 320k -map_metadata 0 -id3v2_version 3 all_voices_final.mp3

ffmpeg \
-i all_voices_final.mp3 \
-i long_background_music.mp3 \
-filter_complex "[1:a]volume='if(lt(t,2),0.1,0.01)':eval=frame[a1]; [0:a][a1]amix=inputs=2:duration=first" \
-c:a libmp3lame output_episode.mp3