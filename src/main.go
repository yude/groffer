package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/google/uuid"
)

//go:embed static/*
var embedDirStatic embed.FS

func main() {
	app := fiber.New()

	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
	}))

	app.Post("/api/render", func(c *fiber.Ctx) error {
		id := uuid.New()
		c.Accepts("text/plain")
		f, err := os.Create(
			fmt.Sprintf("/tmp/%s.roff", id),
		)
		if err != nil {
			fmt.Println(err)
			c.Status(fiber.StatusInternalServerError).JSON(
				ResError{
					Error: "failed to process your request",
				},
			)
		}
		defer f.Close()
		_, err = f.Write(c.BodyRaw())
		if err != nil {
			fmt.Println(err)
			c.Status(fiber.StatusInternalServerError).JSON(
				ResError{
					Error: "failed to write your request into our side",
				},
			)
		}

		cmd := exec.Command(
			"do-render",
			fmt.Sprintf("/tmp/%s.roff", id),
		)
		err = cmd.Run()
		if err != nil || cmd.ProcessState.ExitCode() != 0 {
			fmt.Println(err)
			c.Status(fiber.StatusInternalServerError).JSON(
				ResError{
					Error: "failed to render pdf",
				},
			)
		}

		if err := c.SendFile(
			fmt.Sprintf("/tmp/%s.pdf", id),
		); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(
				ResError{
					Error: "failed to send you the pdf.",
				},
			)
		}

		go func() {
			extensions := [...]string{"pdf", "ps", "roff"}

			for _, ext := range extensions {
				if err := os.Remove(
					fmt.Sprintf("/tmp/%s.%s", id, ext),
				); err != nil {
					fmt.Printf("Failed to delete temporary file: /tmp/%s.%s", id, ext)
					fmt.Println(err)
				}
			}
		}()

		return nil
	})

	log.Fatalln(app.Listen(":3000"))
}
