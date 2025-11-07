package cli

import "github.com/spf13/cobra"

func initFlags(pickCmd *cobra.Command, spotComd *cobra.Command) {
	pickCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	pickCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	pickCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")
	pickCmd.Flags().IntVar(&limit, "limit", 0, "Limit number of films returned (0 = all)")
	pickCmd.Flags().StringVar(&duration, "duration", "long", "Optional duration filter: short, medium, long")

	spotComd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	spotComd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")
	spotComd.Flags().IntVar(&limit, "limit", 0, "Limit number of films returned (0 = all)")
	spotComd.Flags().StringVar(&duration, "duration", "long", "Optional duration filter: short, medium, long")
}
