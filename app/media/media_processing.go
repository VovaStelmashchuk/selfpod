package media

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const TemporaryFolder = "tmp_files"
const AudioFile = "audio.mp3"
const ImageFile = "image.png"

func PrepareNewVideo(audioUrl string, imageUrl string) string {
	err := os.RemoveAll(TemporaryFolder)
	if err != nil {
		log.Printf("Error removing temporary folder: %v", TemporaryFolder)
	}
	err = os.Mkdir(TemporaryFolder, 0755)
	if err != nil {
		log.Printf("Error creating temporary folder: %v", TemporaryFolder)
	}

	metaFiles := downloadFiles(audioUrl, imageUrl)

	return createVideoFile(metaFiles)
}

func createVideoFile(metaFiles MetaFiles) string {
	outputFile := path.Join(TemporaryFolder, "output.mov")

	log.Printf(metaFiles.image)

	log.Printf("Creating video %s by ffmpeg", outputFile)

	cmd := exec.Command(
		"ffmpeg", "-loop", "1", "-i", metaFiles.image, "-i", metaFiles.audio, "-vf", "scale=1280:720", "-tune",
		"stillimage", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-c:a", "aac", "-b:a", "192k", "-shortest", "-loglevel",
		"verbose", outputFile,
	)

	if err := cmd.Run(); err != nil {
		log.Fatalf("FFmpeg error: %v", err)
	}

	log.Printf("Video file %s created successfully", outputFile)

	return outputFile
}

type MetaFiles struct {
	audio string
	image string
}

func downloadFiles(audioUrl string, imageUrl string) MetaFiles {
	audioFile := filepath.Join(TemporaryFolder, AudioFile)
	imageFile := filepath.Join(TemporaryFolder, ImageFile)

	downloadAndSaveFile(audioUrl, audioFile)
	downloadAndSaveFile(imageUrl, imageFile)

	return MetaFiles{
		audio: audioFile,
		image: imageFile,
	}
}
