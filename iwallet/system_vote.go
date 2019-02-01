// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iwallet

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/iost-official/go-iost/rpc/pb"
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:     "vote producerID amount",
	Short:   "Vote a producer",
	Long:    `Vote a producer by given amount of IOSTs`,
	Example: `  iwallet sys vote producer000 1000000 --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			cmd.Usage()
			return fmt.Errorf("please enter the producer ID and the amount")
		}
		_, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			cmd.Usage()
			return fmt.Errorf(`invalid argument "%v" for "amount": %v`, args[1], err)
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "vote", sdk.accountName, args[0], args[1])(cmd, args)
	},
}
var unvoteCmd = &cobra.Command{
	Use:     "unvote producerID amount",
	Short:   "Unvote a producer",
	Long:    `Unvote a producer by given amount of IOSTs`,
	Example: `  iwallet sys unvote producer000 1000000 --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			cmd.Usage()
			return fmt.Errorf("please enter the producer ID and the amount")
		}
		_, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			cmd.Usage()
			return fmt.Errorf(`invalid argument "%v" for "amount": %v`, args[1], err)
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "unvote", sdk.accountName, args[0], args[1])(cmd, args)
	},
}

var location string
var url string
var networkID string
var isPartner bool
var registerCmd = &cobra.Command{
	Use:     "register publicKey",
	Aliases: []string{"reg"},
	Short:   "Register as producer",
	Long:    `Register as producer`,
	Example: `  iwallet sys register XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.Usage()
			return fmt.Errorf("please enter the public key")
		}
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "applyRegister", sdk.accountName, args[0], location, url, networkID, !isPartner)(cmd, args)
	},
}
var unregisterCmd = &cobra.Command{
	Use:     "unregister",
	Aliases: []string{"unreg"},
	Short:   "Unregister from a producer",
	Long:    `Unregister from a producer`,
	Example: `  iwallet sys unregister --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "applyUnregister", sdk.accountName)(cmd, args)
	},
}

var loginCmd = &cobra.Command{
	Use:     "producer-login",
	Aliases: []string{"plogin"},
	Short:   "Producer login as online state",
	Long:    `Producer login as online state`,
	Example: `  iwallet sys plogin --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "logInProducer", sdk.accountName)(cmd, args)
	},
}
var logoutCmd = &cobra.Command{
	Use:     "producer-logout",
	Aliases: []string{"plogout"},
	Short:   "Producer logout as offline state",
	Long:    `Producer logout as offline state`,
	Example: `  iwallet sys plogout --account test0`,
	Args: func(cmd *cobra.Command, args []string) error {
		return checkAccount(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return actionSender("vote_producer.iost", "logOutProducer", sdk.accountName)(cmd, args)
	},
}

var infoCmd = &cobra.Command{
	Use:     "producer-info producerID",
	Aliases: []string{"pinfo"},
	Short:   "Show producer info",
	Long:    `Show producer info`,
	Example: `  iwallet sys pinfo producer000`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.Usage()
			return fmt.Errorf("please enter the producer id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := sdk.GetProducerVoteInfo(&rpcpb.GetProducerVoteInfoRequest{
			Account:        args[0],
			ByLongestChain: sdk.useLongestChain,
		})
		if err != nil {
			return err
		}
		fmt.Println(marshalTextString(info))
		return nil
	},
}

func getProducerList(key string) ([]string, error) {
	response, err := sdk.GetContractStorage(&rpcpb.GetContractStorageRequest{
		Id:  "vote_producer.iost",
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	var list []string
	err = json.Unmarshal([]byte(response.Data), &list)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(list))
	for i, producerKey := range list {
		response, err := sdk.GetContractStorage(&rpcpb.GetContractStorageRequest{
			Id:    "vote_producer.iost",
			Key:   "producerKeyToId",
			Field: producerKey,
		})
		if err != nil {
			return nil, err
		}
		result[i] = response.Data
	}
	return result, nil
}

var listCmd = &cobra.Command{
	Use:     "producer-list",
	Aliases: []string{"plist"},
	Short:   "Show current/pending producer list",
	Long:    `Show current/pending producer list`,
	Example: `  iwallet sys plist`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentList, err := getProducerList("currentProducerList")
		if err != nil {
			return err
		}
		fmt.Println("Current producer list:", currentList)
		pendingList, err := getProducerList("pendingProducerList")
		if err != nil {
			return err
		}
		fmt.Println("Pending producer list:", pendingList)
		return nil
	},
}

func init() {
	systemCmd.AddCommand(voteCmd)
	systemCmd.AddCommand(unvoteCmd)

	systemCmd.AddCommand(registerCmd)
	registerCmd.Flags().StringVarP(&location, "location", "", "", "location info")
	registerCmd.Flags().StringVarP(&url, "url", "", "", "url address")
	registerCmd.Flags().StringVarP(&networkID, "net_id", "", "", "network ID")
	registerCmd.Flags().BoolVarP(&isPartner, "partner", "", false, "if is partner instead of producer")
	systemCmd.AddCommand(unregisterCmd)

	systemCmd.AddCommand(loginCmd)
	systemCmd.AddCommand(logoutCmd)

	systemCmd.AddCommand(infoCmd)
	systemCmd.AddCommand(listCmd)
}