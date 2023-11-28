package media

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const DownloadFolder = "tmp_files"
const AudioFile = "audio.wav"
const ImageFile = "image.jpg"

func PrepareNewVideo(audioUrl string, imageUrl string) {
	os.Mkdir(DownloadFolder, 0755)

	metaFiles := downloadFiles(audioUrl, imageUrl)

	createVideoFile(metaFiles)
}

func createVideoFile(metaFiles MetaFiles) {
	outputFile := path.Join(DownloadFolder, "output.mov")

	println(metaFiles.image)

	cmd := exec.Command("ffmpeg", "-loop", "1", "-i", metaFiles.image, "-i", metaFiles.audio, "-vf", "scale=1920:1080", "-tune", "stillimage", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", "-shortest", "-loglevel", "verbose", outputFile)

	if err := cmd.Run(); err != nil {
		log.Fatalf("FFmpeg error: %v", err)
	}

	log.Printf("Video file %s created successfully", outputFile)
}

type MetaFiles struct {
	audio string
	image string
}

func downloadFiles(audioUrl string, imageUrl string) MetaFiles {
	audioFile := filepath.Join(DownloadFolder, AudioFile)
	imageFile := filepath.Join(DownloadFolder, ImageFile)

	downloadAndSaveFile(audioUrl, audioFile)
	downloadAndSaveFile(imageUrl, imageFile)

	return MetaFiles{
		audio: audioFile,
		image: imageFile,
	}
}
