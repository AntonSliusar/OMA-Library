package handlers

import (
	"context"
	"encoding/json"
	"oma-library/internal/config"
	"oma-library/pkg/storage"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RunMCP(storage *storage.Storage, r2 *storage.R2Client, cfg *config.Config) {
	s := server.NewMCPServer("oma-uploader", "1.0.0")

	tool := mcp.NewTool("list_uploaded_files", mcp.WithDescription("Повертає список файлів які вже завантажені в систему"))
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error){
		files, err := storage.GetAll()
		if err != nil {
			return nil, err
		}
		data, err := json.Marshal(files)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(string(data)), nil
	})


	sseServer := server.NewSSEServer(s, server.WithBaseURL(cfg.MCP.URL))
	sseServer.Start(cfg.MCP.Port)

}