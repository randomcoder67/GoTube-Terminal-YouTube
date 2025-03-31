package download

// This file contains all of the struct definitions used to extract the useful bits of JSON from YouTube

type PLAddRemove struct {
	Header struct {
		PlaylistHeaderRenderer struct {
			SaveButton struct {
				ToggleButtonRenderer struct {
					DefaultServiceEndpoint struct {
						LikeEndpoint struct {
							LikeParams string `json:"likeParams"`
						} `json:"likeEndpoint"`
					} `json:"defaultServiceEndpoint"`
					ToggledServiceEndpoint struct {
						LikeEndpoint struct {
							DislikeParams string `json:"removeLikeParams"`
						} `json:"likeEndpoint"`
					} `json:"toggledServiceEndpoint"`
				} `json:"toggleButtonRenderer"`
			} `json:"saveButton"`
		} `json:"playlistHeaderRenderer"`
	} `json:"header"`
}

type CreatePlaylistResponseJSON struct {
	Error struct {
		Code int `json:"code"`
	} `json:"error"`
	PlaylistId string `json:"playlistId"`
	Actions []struct {
		AddToGuideSectionAction struct {
			Items []struct {
				GuideEntryRenderer struct {
					FormattedTitle struct {
						SimpleText string `json:"simpleText"`
					} `json:"formattedTitle"`
				} `json:"guideEntryRenderer"`
			} `json:"items"`	
		} `json:"addToGuideSectionAction"`
	} `json:"actions"`
}

type ContinuationJSON struct {
	OnResponseReceivedActions []struct {
		AppendContinuationItemsAction struct {
			ContinuationItems []struct {
				PlaylistVideoRenderer struct {
					VideoID   string `json:"videoId"`
					Thumbnail struct {
						Thumbnails []struct {
							URL string `json:"url"`
						} `json:"thumbnails"`
					} `json:"thumbnail"`
					Title struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"title"`
					ShortBylineText struct {
						Runs []struct {
							Text               string `json:"text"`
							NavigationEndpoint struct {
								CommandMetadata struct {
									WebCommandMetadata struct {
										URL string `json:"url"`
									} `json:"webCommandMetadata"`
								} `json:"commandMetadata"`
							} `json:"navigationEndpoint"`
						} `json:"runs"`
					} `json:"shortBylineText"`
					LengthText struct {
						SimpleText string `json:"simpleText"`
					} `json:"lengthText"`
					NavigationEndpoint struct {
						WatchEndpoint struct {
							StartTimeSeconds int `json:"startTimeSeconds"`
						} `json:"watchEndpoint"`
					} `json:"navigationEndpoint"`
					Menu struct {
						MenuRenderer struct {
							Items []struct {
								MenuServiceItemRenderer struct {
									ServiceEndpoint struct {
										PlaylistEditEndpoint struct {
											Params        string `json:"params"`
											ClientActions []struct {
												PlaylistRemoveVideosAction struct {
													SetVideoIds []string `json:"setVideoIds"`
												} `json:"playlistRemoveVideosAction"`
											} `json:"clientActions"`
										} `json:"playlistEditEndpoint"`
									} `json:"serviceEndpoint"`
								} `json:"menuServiceItemRenderer,omitempty"`
							} `json:"items"`
						} `json:"menuRenderer"`
					} `json:"menu"`
					VideoInfo struct {
						Runs []struct {
							Text string `json:"text"`
						} `json:"runs"`
					} `json:"videoInfo"`
				} `json:"playlistVideoRenderer,omitempty"`
				ContinuationItemRenderer struct {
					ContinuationEndpoint struct {
						ContinuationCommand struct {
							Token string `json:"token"`
						} `json:"continuationCommand"`
					} `json:"continuationEndpoint"`
				} `json:"continuationItemRenderer,omitempty"`
			} `json:"continuationItems"`
		} `json:"appendContinuationItemsAction"`
	} `json:"onResponseReceivedActions"`
}

type SubJSON struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						RichGridRenderer struct {
							Contents []struct {
								RichItemRenderer struct {
									Content struct {
										VideoRenderer struct {
											VideoID   string `json:"videoId"`
											Thumbnail struct {
												Thumbnails []struct {
													URL    string `json:"url"`
													Width  int    `json:"width"`
													Height int    `json:"height"`
												} `json:"thumbnails"`
											} `json:"thumbnail"`
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title"`
											DescriptionSnippet struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"descriptionSnippet"`
											PublishedTimeText struct {
												SimpleText string `json:"simpleText"`
											} `json:"publishedTimeText"`
											LengthText struct {
												Accessibility struct {
													AccessibilityData struct {
														Label string `json:"label"`
													} `json:"accessibilityData"`
												} `json:"accessibility"`
												SimpleText string `json:"simpleText"`
											} `json:"lengthText"`
											NavigationEndpoint struct {
												WatchEndpoint struct {
													StartTimeSeconds int `json:"startTimeSeconds"`
												} `json:"watchEndpoint"`
											} `json:"navigationEndpoint"`
											ShortViewCountText struct {
												SimpleText string `json:"simpleText"`
												Runs       []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"shortViewCountText"`
											OwnerText struct {
												Runs []struct {
													Text               string `json:"text"`
													NavigationEndpoint struct {
														CommandMetadata struct {
															WebCommandMetadata struct {
																URL string `json:"url"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
													} `json:"navigationEndpoint"`
												} `json:"runs"`
											} `json:"ownerText"`
										} `json:"videoRenderer"`
									} `json:"content"`
								} `json:"richItemRenderer,omitempty"`
							} `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type HistJSON struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								ContinuationItemRenderer struct {
									ContinuationEndpoint struct {
										ContinuationCommand struct {
											Token string `json:"token"`
										} `json:"continuationCommand"`
									} `json:"continuationEndpoint"`
								} `json:"continuationItemRenderer"`
								ItemSectionRenderer struct {
									Contents []struct {
										VideoRenderer struct {
											VideoID   string `json:"videoId"`
											Thumbnail struct {
												Thumbnails []struct {
													URL    string `json:"url"`
													Width  int    `json:"width"`
													Height int    `json:"height"`
												} `json:"thumbnails"`
											} `json:"thumbnail"`
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title"`
											DescriptionSnippet struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"descriptionSnippet"`
											LengthText struct {
												SimpleText string `json:"simpleText"`
											} `json:"lengthText"`
											NavigationEndpoint struct {
												WatchEndpoint struct {
													StartTimeSeconds int `json:"startTimeSeconds"`
												} `json:"watchEndpoint"`
											} `json:"navigationEndpoint"`
											ShortViewCountText struct {
												SimpleText string `json:"simpleText"`
												Runs       []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"shortViewCountText"`
											OwnerText struct {
												Runs []struct {
													Text               string `json:"text"`
													NavigationEndpoint struct {
														CommandMetadata struct {
															WebCommandMetadata struct {
																URL string `json:"url"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
													} `json:"navigationEndpoint"`
												} `json:"runs"`
											} `json:"ownerText"`
											IsWatched bool `json:"isWatched"`
										} `json:"videoRenderer,omitempty"`
									} `json:"contents"`
								} `json:"itemSectionRenderer,omitempty"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type WLJSON struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								ItemSectionRenderer struct {
									Contents []struct {
										PlaylistVideoListRenderer struct {
											Contents []struct {
												ContinuationItemRenderer struct {
													ContinuationEndpoint struct {
														ContinuationCommand struct {
															Token string `json:"token"`
														} `json:"continuationCommand"`
													} `json:"continuationEndpoint"`
												} `json:"continuationItemRenderer"`
												PlaylistVideoRenderer struct {
													VideoID   string `json:"videoId"`
													Thumbnail struct {
														Thumbnails []struct {
															URL    string `json:"url"`
															Width  int    `json:"width"`
															Height int    `json:"height"`
														} `json:"thumbnails"`
													} `json:"thumbnail"`
													Title struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
														Accessibility struct {
															AccessibilityData struct {
																Label string `json:"label"`
															} `json:"accessibilityData"`
														} `json:"accessibility"`
													} `json:"title"`
													Index struct {
														SimpleText string `json:"simpleText"`
													} `json:"index"`
													ShortBylineText struct {
														Runs []struct {
															Text               string `json:"text"`
															NavigationEndpoint struct {
																CommandMetadata struct {
																	WebCommandMetadata struct {
																		URL string `json:"url"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
															} `json:"navigationEndpoint"`
														} `json:"runs"`
													} `json:"shortBylineText"`
													LengthText struct {
														SimpleText string `json:"simpleText"`
													} `json:"lengthText"`
													NavigationEndpoint struct {
														WatchEndpoint struct {
															StartTimeSeconds int `json:"startTimeSeconds"`
														} `json:"watchEndpoint"`
													} `json:"navigationEndpoint"`
													SetVideoID    string `json:"setVideoId"`
													LengthSeconds string `json:"lengthSeconds"`
													Menu          struct {
														MenuRenderer struct {
															Items []struct {
																MenuServiceItemRenderer struct {
																	ServiceEndpoint struct {
																		PlaylistEditEndpoint struct {
																			Params        string `json:"params"`
																			ClientActions []struct {
																				PlaylistRemoveVideosAction struct {
																					SetVideoIds []string `json:"setVideoIds"`
																				} `json:"playlistRemoveVideosAction"`
																			} `json:"clientActions"`
																		} `json:"playlistEditEndpoint"`
																	} `json:"serviceEndpoint"`
																} `json:"menuServiceItemRenderer"`
															} `json:"items"`
														} `json:"menuRenderer"`
													} `json:"menu"`
													VideoInfo struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"videoInfo"`
												} `json:"playlistVideoRenderer,omitempty"`
											} `json:"contents"`
										} `json:"playlistVideoListRenderer,omitempty"`
									} `json:"contents"`
								} `json:"itemSectionRenderer,omitempty"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type LibraryJSON struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						RichGridRenderer struct {
							Contents []struct {
								RichSectionRenderer struct {
									Content struct {
										RichShelfRenderer struct {
											Subtitle struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"subtitle"`
											
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title"`
											
											Contents []struct {
												RichItemRenderer struct {
													Content struct {

														VideoRenderer struct {
															Thumbnail struct {
																Thumbnails []struct {
																	URL string `json:"url"`
																} `json:"thumbnails"`
															} `json:"thumbnail"`
														} `json:"videoRenderer"`
													
														LockupViewModel struct {
														
															ContentImage struct {
																CollectionThumbnailViewModel struct {
																	PrimaryThumbnail struct {
																		ThumbnailViewModel struct {
																			Image struct {
																				Sources []struct {
																					URL string `json:"url"`
																				} `json:"sources"`
																			} `json:"image"`
																			Overlays []struct {
																				ThumbnailOverlayBadgeViewModel struct {
																					ThumbnailBadges []struct {
																						ThumbnailBadgeViewModel struct {
																							Text string `json:"text"`
																						} `json:"thumbnailBadgeViewModel"`
																					} `json:"thumbnailBadges"`
																				} `json:"thumbnailOverlayBadgeViewModel"`
																			} `json:"overlays"`
																		} `json:"thumbnailViewModel"`
																	} `json:"primaryThumbnail"`
																} `json:"collectionThumbnailViewModel"`
															} `json:"contentImage"`
																					
															Metadata struct {
																LockupMetadataViewModel struct {
																
																	Title struct {
																		Content string `json:"content"`
																	} `json:"title"`

																	Metadata struct {
																		ContentMetadataViewModel struct {
																			MetadataRows []struct {
																				MetadataParts []struct {
																					Text struct {
																						CommandRuns []struct {
																							OnTap struct {
																								InnertubeCommand struct {
																									CommandMetadata struct {
																										WebCommandMetadata struct {
																											URL string `json:"url"`
																										} `json:"webCommandMetadata"`
																									} `json:"commandMetadata"`
																								} `json:"innertubeCommand"`
																							} `json:"onTap"`


																						} `json:"commandRuns"`
																						Content string `json:"content"`
																					} `json:"text"`
																				} `json:"metadataParts"`
																			} `json:"metadataRows"`
																		} `json:"contentMetadataViewModel"`
																	} `json:"metadata"`
																	
																} `json:"lockupMetadataViewModel"`
															} `json:"metadata"`
														} `json:"lockupViewModel"`
														
													} `json:"content"`
												} `json:"richItemRenderer"`
											} `json:"contents"`
										} `json:"richShelfRenderer"`
									} `json:"content"`
								} `json:"richSectionRenderer"`
							} `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}
/*
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title,omitempty"`
											TitleAnnotation struct {
												SimpleText string `json:"simpleText"`
											} `json:"titleAnnotation"`
											Content struct {
												HorizontalListRenderer struct {
													Items []struct {
												RichItemRenderer struct {
													Content struct {
														LockupViewModel struct {
															
															
																	Metadata struct {
																		ContentMetadataViewModel struct {
																			MetadataRows []struct {
																				MetadataParts []struct {
																					Text struct {
																						Content string `json:"content"`
																						CommandRuns []struct {
																							OnTap struct {
																								InnertubeCommand struct {
																									CommandMetadata struct {
																										WebCommandMetadata struct {
																											URL string
																										} `json:"webCommandMetadata"`
																									} `json:"commandMetadata"`
																								} `json:"innertubeCommand"`
																							} `json:"onTap"`
																						} `json:"commandRuns"`
																					}  `json:"text"`
																				} `json:"metadataParts"`
																			}  `json:"metadataRows"`
																		}  `json:"contentMetadataViewModel"`
																	} `json:"metadata"`
																} `json:"lockupMetadataViewModel"`
															} `json:"Metadata"`

														} `json:"lockupViewModel"`


														
														GridVideoRenderer struct {
															Thumbnail struct {
																Thumbnails []struct {
																	URL string `json:"url"`
																	Width int `json:"width"`
																	Height int `json:"height"`
																} `json:"thumbnails"`
															} `json:"thumbnail"`
														} `json:"gridVideoRenderer"`
														GridPlaylistRenderer struct {
															PlaylistID string `json:"playlistId"`
															Thumbnail  struct {
																Thumbnails []struct {
																	URL    string `json:"url"`
																	Width  int    `json:"width"`
																	Height int    `json:"height"`
																} `json:"thumbnails"`
																SampledThumbnailColor struct {
																	Red   int `json:"red"`
																	Green int `json:"green"`
																	Blue  int `json:"blue"`
																} `json:"sampledThumbnailColor"`
																DarkColorPalette struct {
																	Section2Color     int `json:"section2Color"`
																	IconInactiveColor int `json:"iconInactiveColor"`
																	IconDisabledColor int `json:"iconDisabledColor"`
																} `json:"darkColorPalette"`
																VibrantColorPalette struct {
																	IconInactiveColor int `json:"iconInactiveColor"`
																} `json:"vibrantColorPalette"`
															} `json:"thumbnail"`
															Title struct {
																SimpleText string `json:"simpleText"`
															} `json:"title"`
															ShortBylineText struct {
																Runs []struct {
																	Text               string `json:"text"`
																	NavigationEndpoint struct {
																		ClickTrackingParams string `json:"clickTrackingParams"`
																	} `json:"navigationEndpoint"`
																} `json:"runs"`
															} `json:"shortBylineText"`
															VideoCountText struct {
																Runs []struct {
																	Text string `json:"text"`
																} `json:"runs"`
															} `json:"videoCountText"`
															NavigationEndpoint struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																CommandMetadata     struct {
																	WebCommandMetadata struct {
																		URL         string `json:"url"`
																		WebPageType string `json:"webPageType"`
																		RootVe      int    `json:"rootVe"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
																WatchEndpoint struct {
																	VideoID        string `json:"videoId"`
																	PlaylistID     string `json:"playlistId"`
																	Params         string `json:"params"`
																	PlayerParams   string `json:"playerParams"`
																	LoggingContext struct {
																		VssLoggingContext struct {
																			SerializedContextData string `json:"serializedContextData"`
																		} `json:"vssLoggingContext"`
																	} `json:"loggingContext"`
																	WatchEndpointSupportedOnesieConfig struct {
																		HTML5PlaybackOnesieConfig struct {
																			CommonConfig struct {
																				URL string `json:"url"`
																			} `json:"commonConfig"`
																		} `json:"html5PlaybackOnesieConfig"`
																	} `json:"watchEndpointSupportedOnesieConfig"`
																} `json:"watchEndpoint"`
															} `json:"navigationEndpoint"`
															PublishedTimeText struct {
																SimpleText string `json:"simpleText"`
															} `json:"publishedTimeText"`
															VideoCountShortText struct {
																SimpleText string `json:"simpleText"`
															} `json:"videoCountShortText"`
															TrackingParams    string `json:"trackingParams"`
															SidebarThumbnails []struct {
																Thumbnails []struct {
																	URL    string `json:"url"`
																	Width  int    `json:"width"`
																	Height int    `json:"height"`
																} `json:"thumbnails"`
															} `json:"sidebarThumbnails"`
															ThumbnailText struct {
																Runs []struct {
																	Text string `json:"text"`
																	Bold bool   `json:"bold,omitempty"`
																} `json:"runs"`
															} `json:"thumbnailText"`
															LongBylineText struct {
																Runs []struct {
																	Text string `json:"text"`
																} `json:"runs"`
															} `json:"longBylineText"`
															ViewPlaylistText struct {
																Runs []struct {
																	Text               string `json:"text"`
																	NavigationEndpoint struct {
																		ClickTrackingParams string `json:"clickTrackingParams"`
																		CommandMetadata     struct {
																			WebCommandMetadata struct {
																				URL         string `json:"url"`
																				WebPageType string `json:"webPageType"`
																				RootVe      int    `json:"rootVe"`
																				APIURL      string `json:"apiUrl"`
																			} `json:"webCommandMetadata"`
																		} `json:"commandMetadata"`
																		BrowseEndpoint struct {
																			BrowseID string `json:"browseId"`
																		} `json:"browseEndpoint"`
																	} `json:"navigationEndpoint"`
																} `json:"runs"`
															} `json:"viewPlaylistText"`
														} `json:"richItemRenderer,omitempty"`
														
													} `json:"items"`
												} `json:"richShelfRenderer,omitempty"`
											} `json:"content"`
										} `json:"richSectionRenderer,omitempty"`
									} `json:"contents"`
								} `json:"richItemRenderer,omitempty"`
							} `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}
*/

type SearchJSON struct {
	Contents struct {
		TwoColumnSearchResultsRenderer struct {
			PrimaryContents struct {
				SectionListRenderer struct {
					Contents []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								// Playlist
								LockupViewModel struct {
									ContentImage struct {
										CollectionThumbnailViewModel struct {
											PrimaryThumbnail struct {
												ThumbnailViewModel struct {
													Image struct {
														Sources []struct {
															URL    string `json:"url"`
															Width  int    `json:"width"`
															Height int    `json:"height"`
														} `json:"sources"`
													} `json:"image"`
													Overlays []struct {
														ThumbnailOverlayBadgeViewModel struct {
															ThumbnailBadges []struct {
																ThumbnailBadgeViewModel struct {
																	Icon struct {
																		Sources []struct {
																			ClientResource struct {
																				ImageName string `json:"imageName"`
																			} `json:"clientResource"`
																		} `json:"sources"`
																	} `json:"icon"`
																	Text       string `json:"text"`
																	BadgeStyle string `json:"badgeStyle"`
																} `json:"thumbnailBadgeViewModel"`
															} `json:"thumbnailBadges"`
															Position string `json:"position"`
														} `json:"thumbnailOverlayBadgeViewModel"`
													} `json:"overlays"`
												} `json:"thumbnailViewModel"`
											} `json:"primaryThumbnail"`
										} `json:"collectionThumbnailViewModel"`
									} `json:"contentImage"`
									Metadata struct {
										LockupMetadataViewModel struct {
											Title struct {
												Content string `json:"content"`
											} `json:"title"`
											Metadata struct {
												ContentMetadataViewModel struct {
													MetadataRows []struct {
														MetadataParts []struct {
															Text struct {
																Content     string `json:"content"`
																CommandRuns []struct {
																	StartIndex int `json:"startIndex"`
																	Length     int `json:"length"`
																	OnTap      struct {
																		InnertubeCommand struct {
																			ClickTrackingParams string `json:"clickTrackingParams"`
																			CommandMetadata     struct {
																				WebCommandMetadata struct {
																					URL string `json:"url"`
																				} `json:"webCommandMetadata"`
																			} `json:"commandMetadata"`
																		} `json:"innertubeCommand"`
																	} `json:"onTap"`
																} `json:"commandRuns"`
															} `json:"text"`
														} `json:"metadataParts"`
													} `json:"metadataRows"`
													Delimiter string `json:"delimiter"`
												} `json:"contentMetadataViewModel"`
											} `json:"metadata"`
										} `json:"lockupMetadataViewModel"`
									} `json:"metadata"`
								} `json:"lockupViewModel"`


								/*
								PlaylistRenderer struct {
									PlaylistID string `json:"playlistId"`
									Thumbnails []struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"thumbnails"`
									Title struct {
										SimpleText string `json:"simpleText"`
									} `json:"title"`
									ShortBylineText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"shortBylineText"`
									VideoCountText struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"videoCountText"`
									PublishedTimeText struct {
										SimpleText string `json:"simpleText"`
									} `json:"publishedTimeText"`
									VideoCountShortText struct {
										SimpleText string `json:"simpleText"`
									} `json:"videoCountShortText"`
									TrackingParams    string `json:"trackingParams"`
									SidebarThumbnails []struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"sidebarThumbnails"`
									ThumbnailText struct {
										Runs []struct {
											Text string `json:"text"`
											Bold bool   `json:"bold,omitempty"`
										} `json:"runs"`
									} `json:"thumbnailText"`
									LongBylineText struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"longBylineText"`
								} `json:"playlistRenderer,omitempty"`
								*/
								// Video
								VideoRenderer struct {
									VideoID   string `json:"videoId"`
									Thumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"thumbnail"`
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									PublishedTimeText struct {
										SimpleText string `json:"simpleText"`
									} `json:"publishedTimeText"`
									LengthText struct {
										SimpleText string `json:"simpleText"`
									} `json:"lengthText"`
									ViewCountText struct {
										SimpleText string `json:"simpleText"`
										Runs       []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"viewCountText"`
									NavigationEndpoint struct {
										WatchEndpoint struct {
											StartTimeSeconds int `json:"startTimeSeconds"`
										} `json:"watchEndpoint"`
									} `json:"navigationEndpoint"`
									ShortViewCountText struct {
										SimpleText string `json:"simpleText"`
										Runs       []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"shortViewCountText"`
									OwnerText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												CommandMetadata struct {
													WebCommandMetadata struct {
														URL string `json:"url"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"ownerText"`
									IsWatched              bool `json:"isWatched"`
									InlinePlaybackEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										WatchEndpoint struct {
											VideoID              string `json:"videoId"`
											Params               string `json:"params"`
											PlayerParams         string `json:"playerParams"`
											PlayerExtraURLParams []struct {
												Key   string `json:"key"`
												Value string `json:"value"`
											} `json:"playerExtraUrlParams"`
											WatchEndpointSupportedOnesieConfig struct {
												HTML5PlaybackOnesieConfig struct {
													CommonConfig struct {
														URL string `json:"url"`
													} `json:"commonConfig"`
												} `json:"html5PlaybackOnesieConfig"`
											} `json:"watchEndpointSupportedOnesieConfig"`
										} `json:"watchEndpoint"`
									} `json:"inlinePlaybackEndpoint"`
								} `json:"videoRenderer,omitempty"`
							} `json:"contents"`
						} `json:"itemSectionRenderer"`
					} `json:"contents"`
				} `json:"sectionListRenderer"`
			} `json:"primaryContents"`
		} `json:"twoColumnSearchResultsRenderer"`
	} `json:"contents"`
}

type RecJSON struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						RichGridRenderer struct {
							Contents []struct {
								RichItemRenderer struct {
									Content struct {
										VideoRenderer struct {
											VideoID   string `json:"videoId"`
											Thumbnail struct {
												Thumbnails []struct {
													URL    string `json:"url"`
													Width  int    `json:"width"`
													Height int    `json:"height"`
												} `json:"thumbnails"`
											} `json:"thumbnail"`
											Title struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"title"`
											DescriptionSnippet struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"descriptionSnippet"`
											PublishedTimeText struct {
												SimpleText string `json:"simpleText"`
											} `json:"publishedTimeText"`
											LengthText struct {
												Accessibility struct {
													AccessibilityData struct {
														Label string `json:"label"`
													} `json:"accessibilityData"`
												} `json:"accessibility"`
												SimpleText string `json:"simpleText"`
											} `json:"lengthText"`
											NavigationEndpoint struct {
												WatchEndpoint struct {
													StartTimeSeconds int `json:"startTimeSeconds"`
												} `json:"watchEndpoint"`
											} `json:"navigationEndpoint"`
											ShortViewCountText struct {
												SimpleText string `json:"simpleText"`
												Runs       []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"shortViewCountText"`
											OwnerText struct {
												Runs []struct {
													Text               string `json:"text"`
													NavigationEndpoint struct {
														CommandMetadata struct {
															WebCommandMetadata struct {
																URL string `json:"url"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
													} `json:"navigationEndpoint"`
												} `json:"runs"`
											} `json:"ownerText"`
										} `json:"videoRenderer"`
									} `json:"content"`
								} `json:"richItemRenderer,omitempty"`
							} `json:"contents"`
						} `json:"richGridRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

type DLResponse struct {
	StreamingData struct {
		Formats []struct {
			Itag     int    `json:"itag"`
			URL      string `json:"url"`
			MimeType string `json:"mimeType"`
			Bitrate  int    `json:"bitrate"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FPS      int    `json:"fps"`
		} `json:"formats"`
		AdaptiveFormats []struct {
			Itag     int    `json:"itag"`
			URL      string `json:"url"`
			MimeType string `json:"mimeType"`
			Bitrate  int    `json:"bitrate"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FPS      int    `json:"fps"`
		} `json:"adaptiveFormats"`
		DashManifestURL string `json:"dashManifestUrl"`
		HLSManifestURL  string `json:"hlsManifestUrl"`
	} `json:"streamingData"`
}

type PlaylistAddDataJSON struct {
	Contents []struct {
		AddToPlaylistRenderer struct {
			Playlists []struct {
				PlaylistAddToOptionRenderer struct {
					PlaylistID string `json:"playlistId"`
					Title      struct {
						SimpleText string `json:"simpleText"`
					} `json:"title"`
					Privacy                string `json:"privacy"`
					ContainsSelectedVideos string `json:"containsSelectedVideos"`
				} `json:"PlaylistAddToOptionRenderer,omitempty"`
			} `json:"playlists"`
		} `json:"addToPlaylistRenderer"`
	} `json:"contents"`
}

type VidPagePlayerResp struct {
	/*
		StreamingData struct {
			ExpiresInSeconds string `json:"expiresInSeconds"`
			Formats		  []struct {
				Itag			 int	`json:"itag"`
				URL			  string `json:"url"`
				MimeType		 string `json:"mimeType"`
				Bitrate		  int	`json:"bitrate"`
				Width			int	`json:"width"`
				Height		   int	`json:"height"`
				LastModified	 string `json:"lastModified"`
				ContentLength	string `json:"contentLength,omitempty"`
				Quality		  string `json:"quality"`
				Fps			  int	`json:"fps"`
				QualityLabel	 string `json:"qualityLabel"`
				ProjectionType   string `json:"projectionType"`
				AverageBitrate   int	`json:"averageBitrate,omitempty"`
				AudioQuality	 string `json:"audioQuality"`
				ApproxDurationMs string `json:"approxDurationMs"`
				AudioSampleRate  string `json:"audioSampleRate"`
				AudioChannels	int	`json:"audioChannels"`
			} `json:"formats"`
			AdaptiveFormats []struct {
				Itag	  int	`json:"itag"`
				URL	   string `json:"url"`
				MimeType  string `json:"mimeType"`
				Bitrate   int	`json:"bitrate"`
				Width	 int	`json:"width,omitempty"`
				Height	int	`json:"height,omitempty"`
				InitRange struct {
					Start string `json:"start"`
					End   string `json:"end"`
				} `json:"initRange"`
				IndexRange struct {
					Start string `json:"start"`
					End   string `json:"end"`
				} `json:"indexRange"`
				LastModified	 string `json:"lastModified"`
				ContentLength	string `json:"contentLength"`
				Quality		  string `json:"quality"`
				Fps			  int	`json:"fps,omitempty"`
				QualityLabel	 string `json:"qualityLabel,omitempty"`
				ProjectionType   string `json:"projectionType"`
				AverageBitrate   int	`json:"averageBitrate"`
				ApproxDurationMs string `json:"approxDurationMs"`
				ColorInfo		struct {
					Primaries			   string `json:"primaries"`
					TransferCharacteristics string `json:"transferCharacteristics"`
					MatrixCoefficients	  string `json:"matrixCoefficients"`
				} `json:"colorInfo,omitempty"`
				HighReplication bool	`json:"highReplication,omitempty"`
				AudioQuality	string  `json:"audioQuality,omitempty"`
				AudioSampleRate string  `json:"audioSampleRate,omitempty"`
				AudioChannels   int	 `json:"audioChannels,omitempty"`
				LoudnessDb	  float64 `json:"loudnessDb,omitempty"`
			} `json:"adaptiveFormats"`
		} `json:"streamingData"`
		Captions struct {
			PlayerCaptionsTracklistRenderer struct {
				CaptionTracks []struct {
					BaseURL string `json:"baseUrl"`
					Name	struct {
						SimpleText string `json:"simpleText"`
					} `json:"name"`
					VssID		  string `json:"vssId"`
					LanguageCode   string `json:"languageCode"`
					Kind		   string `json:"kind"`
					IsTranslatable bool   `json:"isTranslatable"`
					TrackName	  string `json:"trackName"`
				} `json:"captionTracks"`
				AudioTracks []struct {
					CaptionTrackIndices []int `json:"captionTrackIndices"`
				} `json:"audioTracks"`
				TranslationLanguages []struct {
					LanguageCode string `json:"languageCode"`
					LanguageName struct {
						SimpleText string `json:"simpleText"`
					} `json:"languageName"`
				} `json:"translationLanguages"`
				DefaultAudioTrackIndex int `json:"defaultAudioTrackIndex"`
			} `json:"playerCaptionsTracklistRenderer"`
		} `json:"captions"`
	*/
	PlaybackTracking struct {
		VideoStatsPlaybackURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videoStatsPlaybackURL"`
		VideoStatsWatchtimeURL struct {
			BaseURL string `json:"baseUrl"`
		} `json:"videoStatsWatchtimeUrl"`
	} `json:"playbackTracking"`
	VideoDetails struct {
		VideoID          string `json:"videoId"`
		Title            string `json:"title"`
		LengthSeconds    string `json:"lengthSeconds"`
		ChannelID        string `json:"channelId"`
		ShortDescription string `json:"shortDescription"`
		Thumbnail        struct {
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
		} `json:"thumbnail"`
		AllowRatings  bool   `json:"allowRatings"`
		ViewCount     string `json:"viewCount"`
		Author        string `json:"author"`
		IsLiveContent bool   `json:"isLiveContent"`
	} `json:"videoDetails"`
	Annotations []struct {
		PlayerAnnotationsExpandedRenderer struct {
			FeaturedChannel struct {
				Watermark struct {
					Thumbnails []struct {
						URL    string `json:"url"`
						Width  int    `json:"width"`
						Height int    `json:"height"`
					} `json:"thumbnails"`
				} `json:"watermark"`
				ChannelName     string `json:"channelName"`
				SubscribeButton struct {
					SubscribeButtonRenderer struct {
						ButtonText struct {
							Runs []struct {
								Text string `json:"text"`
							} `json:"runs"`
						} `json:"buttonText"`
						Subscribed       bool   `json:"subscribed"`
						ChannelID        string `json:"channelId"`
						ServiceEndpoints []struct {
							CommandMetadata struct {
								WebCommandMetadata struct {
									SendPost bool   `json:"sendPost"`
									APIURL   string `json:"apiUrl"`
								} `json:"webCommandMetadata"`
							} `json:"commandMetadata"`
							SubscribeEndpoint struct {
								ChannelIds []string `json:"channelIds"`
								Params     string   `json:"params"`
							} `json:"subscribeEndpoint,omitempty"`
							UnsubscribeEndpoint struct {
								ChannelIds []string `json:"channelIds"`
								Params     string   `json:"params"`
							} `json:"unsubscribeEndpoint"`
							SignalServiceEndpoint struct {
								Signal  string `json:"signal"`
								Actions []struct {
									ClickTrackingParams string `json:"clickTrackingParams"`
									OpenPopupAction     struct {
										Popup struct {
											ConfirmDialogRenderer struct {
												TrackingParams string `json:"trackingParams"`
												DialogMessages []struct {
													Runs []struct {
														Text string `json:"text"`
													} `json:"runs"`
												} `json:"dialogMessages"`
												ConfirmButton struct {
													ButtonRenderer struct {
														Style      string `json:"style"`
														Size       string `json:"size"`
														IsDisabled bool   `json:"isDisabled"`
														Text       struct {
															Runs []struct {
																Text string `json:"text"`
															} `json:"runs"`
														} `json:"text"`
														ServiceEndpoint struct {
															UnsubscribeEndpoint struct {
																ChannelIds []string `json:"channelIds"`
																Params     string   `json:"params"`
															} `json:"unsubscribeEndpoint"`
															SubscribeEndpoint struct {
																ChannelIds []string `json:"channelIds"`
																Params     string   `json:"params"`
															} `json:"subscribeEndpoint,omitempty"`
														} `json:"serviceEndpoint"`
													} `json:"buttonRenderer"`
												} `json:"confirmButton"`
											} `json:"confirmDialogRenderer"`
										} `json:"popup"`
									} `json:"openPopupAction"`
								} `json:"actions"`
							} `json:"signalServiceEndpoint,omitempty"`
						} `json:"serviceEndpoints"`
					} `json:"subscribeButtonRenderer"`
				} `json:"subscribeButton"`
			} `json:"featuredChannel"`
		} `json:"playerAnnotationsExpandedRenderer"`
	} `json:"annotations"`
	Endscreen struct {
		EndscreenRenderer struct {
			Elements []struct {
				EndscreenElementRenderer struct {
					HovercardButton struct {
						SubscribeButtonRenderer struct {
							Subscribed       bool
							ServiceEndpoints []struct {
								UnsubscribeEndpoint struct {
									ChannelIds []string `json:"channelIds"`
									Params     string   `json:"params"`
								} `json:"unsubscribeEndpoint"`
								SubscribeEndpoint struct {
									ChannelIds []string `json:"channelIds"`
									Params     string   `json:"params"`
								} `json:"subscribeEndpoint,omitempty"`
								SignalServiceEndpoint struct {
									Actions []struct {
										OpenPopupAction struct {
											Popup struct {
												ConfirmDialogRenderer struct {
													TrackingParams string `json:"trackingParams"`
													DialogMessages []struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"dialogMessages"`
													ConfirmButton struct {
														ButtonRenderer struct {
															Style      string `json:"style"`
															Size       string `json:"size"`
															IsDisabled bool   `json:"isDisabled"`
															Text       struct {
																Runs []struct {
																	Text string `json:"text"`
																} `json:"runs"`
															} `json:"text"`
															ServiceEndpoint struct {
																UnsubscribeEndpoint struct {
																	ChannelIds []string `json:"channelIds"`
																	Params     string   `json:"params"`
																} `json:"unsubscribeEndpoint"`
																SubscribeEndpoint struct {
																	ChannelIds []string `json:"channelIds"`
																	Params     string   `json:"params"`
																} `json:"subscribeEndpoint,omitempty"`
															} `json:"serviceEndpoint"`
														} `json:"buttonRenderer"`
													} `json:"confirmButton"`
												} `json:"confirmDialogRenderer"`
											} `json:"popup"`
										} `json:"openPopupAction"`
									} `json:"actions"`
								} `json:"signalServiceEndpoint,omitempty"`
							} `json:"serviceEndpoints"`
						} `json:"subscribeButtonRenderer"`
					} `json:"hovercardButton"`
				} `json:"endscreenElementRenderer"`
			} `json:"elements"`
		} `json:"endscreenRenderer"`
	} `json:"endscreen"`
	PlayerConfig struct {
		WebPlayerConfig struct {
			WebPlayerActionsPorting struct {
				SubscribeCommand struct {
					SubscribeEndpoint struct {
						Params string `json:"params"`
					} `json:"subscribeEndpoint"`
				} `json:"subscribeCommand"`
				UnsubscribeCommand struct {
					UnsubscribeEndpoint struct {
						Params string `json:"params"`
					} `json:"unsubscribeEndpoint"`
				} `json:"unsubscribeCommand"`
			} `json:"webPlayerActionsPorting"`
		} `json:"webPlayerConfig"`
		PlaybackStartConfig struct {
			StartSeconds int `json:"startSeconds"`
		} `json:"playbackStartConfig"`
	} `json:"playerConfig"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Category         string `json:"category"`
			PublishDate      string `json:"publishDate"`
			OwnerChannelName string `json:"ownerChannelName"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
	Attestation struct {
		PlayerAttestationRenderer struct {
			Challenge    string `json:"challenge"`
			BotguardData struct {
				Program            string `json:"program"`
				InterpreterSafeURL struct {
					PrivateDoNotAccessOrElseTrustedResourceURLWrappedValue string `json:"privateDoNotAccessOrElseTrustedResourceUrlWrappedValue"`
				} `json:"interpreterSafeUrl"`
				ServerEnvironment int `json:"serverEnvironment"`
			} `json:"botguardData"`
		} `json:"playerAttestationRenderer"`
	} `json:"attestation"`
}

type VidPageInitialData struct {
	Contents struct {
		TwoColumnWatchNextResults struct {
			Results struct {
				Results struct {
					Contents []struct {
						// This is about the main video
						VideoPrimaryInfoRenderer struct {
							ViewCount struct {
								VideoViewCountRenderer struct {
									ViewCount struct {
										SimpleText string `json:"simpleText"`
									} `json:"viewCount"`
									ShortViewCount struct {
										SimpleText string `json:"simpleText"`
									} `json:"shortViewCount"`
								} `json:"videoViewCountRenderer"`
							} `json:"viewCount"`
							DateText struct {
								SimpleText string `json:"simpleText"`
							} `json:"dateText"`
							RelativeDateText struct {
								SimpleText string `json:"simpleText"`
							} `json:"relativeDateText"`
							VideoActions struct {
								MenuRenderer struct {
									TopLevelButtons []struct {
										SegmentedLikeDislikeButtonViewModel struct {
											LikeButtonViewModel struct {
												LikeButtonViewModel struct {
													ToggleButtonViewModel struct {
														ToggleButtonViewModel struct {
															DefaultButtonViewModel struct {
																ButtonViewModel struct {
																	IconName string `json:"iconName"`
																	Title    string `json:"title"`
																	OnTap    struct {
																		SerialCommand struct {
																			Commands []struct {
																				InnertubeCommand struct {
																					ModalEndpoint struct {
																						Modal struct {
																							ModalWithTitleAndButtonRenderer struct {
																								Button struct {
																									ButtonRenderer struct {
																										NavigationEndpoint struct {
																											SignInEndpoint struct {
																												NextEndpoint struct {
																													CommandMetadata struct {
																														WebCommandMetadata struct {
																															ApiURL string `json:"apiUrl"`
																														} `json:"webCommandMetadata"`
																													} `json:"commandMetadata"`
																													LikeEndpoint struct {
																														Status string `json:"status"`
																														Target struct {
																															VideoID string `json:"videoId"`
																														} `json:"target"`
																														LikeParams string `json:"likeParams"`
																													} `json:"likeEndpoint"`
																												} `json:"nextEndpoint"`
																											} `json:"signInEndpoint"`
																										} `json:"navigationEndpoint"`
																									} `json:"buttonRenderer"`
																								} `json:"button"`
																							} `json:"modalWithTitleAndButtonRenderer"`
																						} `json:"modal"`
																					} `json:"modalEndpoint"`
																				} `json:"innertubeCommand,omitempty"`
																			} `json:"commands"`
																		} `json:"serialCommand"`
																	} `json:"onTap"`
																} `json:"buttonViewModel"`
															} `json:"defaultButtonViewModel"`
															ToggledButtonViewModel struct {
																ButtonViewModel struct {
																	IconName string `json:"iconName"`
																	Title    string `json:"title"`
																	OnTap    struct {
																		SerialCommand struct {
																			Commands []struct {
																				InnertubeCommand struct {
																					CommandMetadata struct {
																						WebCommandMetadata struct {
																							ApiURL string `json:"apiUrl"`
																						} `json:"webCommandMetadata"`
																					} `json:"commandMetadata"`
																					LikeEndpoint struct {
																						Status string `json:"status"`
																						Target struct {
																							VideoID string `json:"videoId"`
																						} `json:"target"`
																						RemoveLikeParams string `json:"removeLikeParams"`
																					} `json:"likeEndpoint"`
																				} `json:"innertubeCommand,omitempty"`
																			} `json:"commands"`
																		} `json:"serialCommand"`
																	} `json:"onTap"`
																} `json:"buttonViewModel"`
															} `json:"toggledButtonViewModel"`
														} `json:"toggleButtonViewModel"`
													} `json:"toggleButtonViewModel"`
													LikeStatusEntity struct {
														Key        string `json:"key"`
														LikeStatus string `json:"likeStatus"`
													} `json:"likeStatusEntity"`
												} `json:"likeButtonViewModel"`
											} `json:"likeButtonViewModel"`
											DislikeButtonViewModel struct {
												DislikeButtonViewModel struct {
													ToggleButtonViewModel struct {
														ToggleButtonViewModel struct {
															DefaultButtonViewModel struct {
																ButtonViewModel struct {
																	IconName string `json:"iconName"`
																	Title    string `json:"title"`
																	OnTap    struct {
																		SerialCommand struct {
																			Commands []struct {
																				InnertubeCommand struct {
																					ModalEndpoint struct {
																						Modal struct {
																							ModalWithTitleAndButtonRenderer struct {
																								Button struct {
																									ButtonRenderer struct {
																										NavigationEndpoint struct {
																											SignInEndpoint struct {
																												NextEndpoint struct {
																													CommandMetadata struct {
																														WebCommandMetadata struct {
																															ApiURL string `json:"apiUrl"`
																														} `json:"webCommandMetadata"`
																													} `json:"commandMetadata"`
																													LikeEndpoint struct {
																														Status string `json:"status"`
																														Target struct {
																															VideoID string `json:"videoId"`
																														} `json:"target"`
																														DislikeParams string `json:"dislikeParams"`
																													} `json:"likeEndpoint"`
																												} `json:"nextEndpoint"`
																											} `json:"signInEndpoint"`
																										} `json:"navigationEndpoint"`
																									} `json:"buttonRenderer"`
																								} `json:"button"`
																							} `json:"modalWithTitleAndButtonRenderer"`
																						} `json:"modal"`
																					} `json:"modalEndpoint"`
																				} `json:"innertubeCommand,omitempty"`
																			} `json:"commands"`
																		} `json:"serialCommand"`
																	} `json:"onTap"`
																} `json:"buttonViewModel"`
															} `json:"defaultButtonViewModel"`
															ToggledButtonViewModel struct {
																ButtonViewModel struct {
																	IconName string `json:"iconName"`
																	Title    string `json:"title"`
																	OnTap    struct {
																		SerialCommand struct {
																			Commands []struct {
																				InnertubeCommand struct {
																					CommandMetadata struct {
																						WebCommandMetadata struct {
																							ApiURL string `json:"apiUrl"`
																						} `json:"webCommandMetadata"`
																					} `json:"commandMetadata"`
																					LikeEndpoint struct {
																						Status string `json:"status"`
																						Target struct {
																							VideoID string `json:"videoId"`
																						} `json:"target"`
																						RemoveLikeParams string `json:"removeLikeParams"`
																					} `json:"likeEndpoint"`
																				} `json:"innertubeCommand,omitempty"`
																			} `json:"commands"`
																		} `json:"serialCommand"`
																	} `json:"onTap"`
																} `json:"buttonViewModel"`
															} `json:"toggledButtonViewModel"`
														} `json:"toggleButtonViewModel"`
													} `json:"toggleButtonViewModel"`
												} `json:"dislikeButtonViewModel"`
											} `json:"dislikeButtonViewModel"`
										} `json:"segmentedLikeDislikeButtonViewModel,omitempty"`
									} `json:"topLevelButtons"`
								} `json:"menuRenderer"`
							} `json:"videoActions"`
						} `json:"videoPrimaryInfoRenderer,omitempty"`
						VideoSecondaryInfoRenderer struct {
							Owner struct {
								VideoOwnerRenderer struct {
									Thumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"thumbnail"`
								} `json:"videoOwnerRenderer"`
								SubscriptionButton struct {
									Subscribed bool `json:"subscribed"`
								} `json:"subscriptionButton"`
							} `json:"owner"`
						} `json:"videoSecondaryInfoRenderer"`
					} `json:"contents"`
				} `json:"results"`
			} `json:"results"`

			// Suggestion video renderer
			SecondaryResults struct {
				SecondaryResults struct {
					Results []struct {
						ItemSectionRenderer struct {
							Contents []struct {
								CompactVideoRenderer struct {
									VideoID   string `json:"videoId"`
									Thumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"thumbnail"`
									Title struct {
										Accessibility struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibility"`
										SimpleText string `json:"simpleText"`
									} `json:"title"`
									LongBylineText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														URL         string `json:"url"`
														WebPageType string `json:"webPageType"`
														RootVe      int    `json:"rootVe"`
														APIURL      string `json:"apiUrl"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												BrowseEndpoint struct {
													BrowseID         string `json:"browseId"`
													CanonicalBaseURL string `json:"canonicalBaseUrl"`
												} `json:"browseEndpoint"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"longBylineText"`
									PublishedTimeText struct {
										SimpleText string `json:"simpleText"`
									} `json:"publishedTimeText"`
									ViewCountText struct {
										SimpleText string `json:"simpleText"`
										Runs       []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"viewCountText"`
									LengthText struct {
										Accessibility struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibility"`
										SimpleText string `json:"simpleText"`
									} `json:"lengthText"`
									NavigationEndpoint struct {
										ClickTrackingParams string `json:"clickTrackingParams"`
										CommandMetadata     struct {
											WebCommandMetadata struct {
												URL         string `json:"url"`
												WebPageType string `json:"webPageType"`
												RootVe      int    `json:"rootVe"`
											} `json:"webCommandMetadata"`
										} `json:"commandMetadata"`
										WatchEndpoint struct {
											VideoID                            string `json:"videoId"`
											Params                             string `json:"params"`
											StartTimeSeconds                   int    `json:"startTimeSeconds"`
											Nofollow                           bool   `json:"nofollow"`
											WatchEndpointSupportedOnesieConfig struct {
												HTML5PlaybackOnesieConfig struct {
													CommonConfig struct {
														URL string `json:"url"`
													} `json:"commonConfig"`
												} `json:"html5PlaybackOnesieConfig"`
											} `json:"watchEndpointSupportedOnesieConfig"`
										} `json:"watchEndpoint"`
									} `json:"navigationEndpoint"`
									ShortBylineText struct {
										Runs []struct {
											Text               string `json:"text"`
											NavigationEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														URL         string `json:"url"`
														WebPageType string `json:"webPageType"`
														RootVe      int    `json:"rootVe"`
														APIURL      string `json:"apiUrl"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												BrowseEndpoint struct {
													BrowseID         string `json:"browseId"`
													CanonicalBaseURL string `json:"canonicalBaseUrl"`
												} `json:"browseEndpoint"`
											} `json:"navigationEndpoint"`
										} `json:"runs"`
									} `json:"shortBylineText"`
									Badges []struct {
										MetadataBadgeRenderer struct {
											Style          string `json:"style"`
											Label          string `json:"label"`
											TrackingParams string `json:"trackingParams"`
										} `json:"metadataBadgeRenderer"`
									} `json:"badges"`
									ChannelThumbnail struct {
										Thumbnails []struct {
											URL    string `json:"url"`
											Width  int    `json:"width"`
											Height int    `json:"height"`
										} `json:"thumbnails"`
									} `json:"channelThumbnail"`
									TrackingParams     string `json:"trackingParams"`
									ShortViewCountText struct {
										Accessibility struct {
											AccessibilityData struct {
												Label string `json:"label"`
											} `json:"accessibilityData"`
										} `json:"accessibility"`
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
										SimpleText string `json:"simpleText"`
									} `json:"shortViewCountText"`
									Menu struct {
										MenuRenderer struct {
											Items []struct {
												MenuServiceItemRenderer struct {
													Text struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"text"`
													Icon struct {
														IconType string `json:"iconType"`
													} `json:"icon"`
													ServiceEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																SendPost bool `json:"sendPost"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														SignalServiceEndpoint struct {
															Signal  string `json:"signal"`
															Actions []struct {
																ClickTrackingParams  string `json:"clickTrackingParams"`
																AddToPlaylistCommand struct {
																	OpenMiniplayer      bool   `json:"openMiniplayer"`
																	OpenListPanel       bool   `json:"openListPanel"`
																	VideoID             string `json:"videoId"`
																	ListType            string `json:"listType"`
																	OnCreateListCommand struct {
																		ClickTrackingParams string `json:"clickTrackingParams"`
																		CommandMetadata     struct {
																			WebCommandMetadata struct {
																				SendPost bool   `json:"sendPost"`
																				APIURL   string `json:"apiUrl"`
																			} `json:"webCommandMetadata"`
																		} `json:"commandMetadata"`
																		CreatePlaylistServiceEndpoint struct {
																			VideoIds []string `json:"videoIds"`
																			Params   string   `json:"params"`
																		} `json:"createPlaylistServiceEndpoint"`
																	} `json:"onCreateListCommand"`
																	VideoIds []string `json:"videoIds"`
																} `json:"addToPlaylistCommand,omitempty"`
																OpenPopupAction struct {
																	Popup struct {
																		NotificationActionRenderer struct {
																			ResponseText struct {
																				SimpleText string `json:"simpleText"`
																			} `json:"responseText"`
																			TrackingParams string `json:"trackingParams"`
																		} `json:"notificationActionRenderer"`
																	} `json:"popup"`
																	PopupType string `json:"popupType"`
																} `json:"openPopupAction,omitempty"`
															} `json:"actions"`
														} `json:"signalServiceEndpoint"`
													} `json:"serviceEndpoint"`
													TrackingParams string `json:"trackingParams"`
												} `json:"menuServiceItemRenderer,omitempty"`
												MenuServiceItemRenderer0 struct {
													Text struct {
														Runs []struct {
															Text string `json:"text"`
														} `json:"runs"`
													} `json:"text"`
													Icon struct {
														IconType string `json:"iconType"`
													} `json:"icon"`
													ServiceEndpoint struct {
														ClickTrackingParams string `json:"clickTrackingParams"`
														CommandMetadata     struct {
															WebCommandMetadata struct {
																SendPost bool   `json:"sendPost"`
																APIURL   string `json:"apiUrl"`
															} `json:"webCommandMetadata"`
														} `json:"commandMetadata"`
														ShareEntityServiceEndpoint struct {
															SerializedShareEntity string `json:"serializedShareEntity"`
															Commands              []struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																OpenPopupAction     struct {
																	Popup struct {
																		UnifiedSharePanelRenderer struct {
																			TrackingParams     string `json:"trackingParams"`
																			ShowLoadingSpinner bool   `json:"showLoadingSpinner"`
																		} `json:"unifiedSharePanelRenderer"`
																	} `json:"popup"`
																	PopupType string `json:"popupType"`
																	BeReused  bool   `json:"beReused"`
																} `json:"openPopupAction"`
															} `json:"commands"`
														} `json:"shareEntityServiceEndpoint"`
													} `json:"serviceEndpoint"`
													TrackingParams string `json:"trackingParams"`
													HasSeparator   bool   `json:"hasSeparator"`
												} `json:"menuServiceItemRenderer,omitempty"`
											} `json:"items"`
											TrackingParams string `json:"trackingParams"`
											Accessibility  struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"accessibility"`
											TargetID string `json:"targetId"`
										} `json:"menuRenderer"`
									} `json:"menu"`
									ThumbnailOverlays []struct {
										ThumbnailOverlayResumePlaybackRenderer struct {
											PercentDurationWatched int `json:"percentDurationWatched"`
										} `json:"thumbnailOverlayResumePlaybackRenderer,omitempty"`
										ThumbnailOverlayTimeStatusRenderer struct {
											Text struct {
												Accessibility struct {
													AccessibilityData struct {
														Label string `json:"label"`
													} `json:"accessibilityData"`
												} `json:"accessibility"`
												SimpleText string `json:"simpleText"`
											} `json:"text"`
											Style string `json:"style"`
										} `json:"thumbnailOverlayTimeStatusRenderer,omitempty"`
										ThumbnailOverlayToggleButtonRenderer struct {
											IsToggled     bool `json:"isToggled"`
											UntoggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"untoggledIcon"`
											ToggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"toggledIcon"`
											UntoggledTooltip         string `json:"untoggledTooltip"`
											ToggledTooltip           string `json:"toggledTooltip"`
											UntoggledServiceEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														SendPost bool   `json:"sendPost"`
														APIURL   string `json:"apiUrl"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												PlaylistEditEndpoint struct {
													PlaylistID string `json:"playlistId"`
													Actions    []struct {
														AddedVideoID string `json:"addedVideoId"`
														Action       string `json:"action"`
													} `json:"actions"`
												} `json:"playlistEditEndpoint"`
											} `json:"untoggledServiceEndpoint"`
											ToggledServiceEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														SendPost bool   `json:"sendPost"`
														APIURL   string `json:"apiUrl"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												PlaylistEditEndpoint struct {
													PlaylistID string `json:"playlistId"`
													Actions    []struct {
														Action         string `json:"action"`
														RemovedVideoID string `json:"removedVideoId"`
													} `json:"actions"`
												} `json:"playlistEditEndpoint"`
											} `json:"toggledServiceEndpoint"`
											UntoggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"untoggledAccessibility"`
											ToggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"toggledAccessibility"`
											TrackingParams string `json:"trackingParams"`
										} `json:"thumbnailOverlayToggleButtonRenderer,omitempty"`
										ThumbnailOverlayToggleButtonRenderer0 struct {
											UntoggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"untoggledIcon"`
											ToggledIcon struct {
												IconType string `json:"iconType"`
											} `json:"toggledIcon"`
											UntoggledTooltip         string `json:"untoggledTooltip"`
											ToggledTooltip           string `json:"toggledTooltip"`
											UntoggledServiceEndpoint struct {
												ClickTrackingParams string `json:"clickTrackingParams"`
												CommandMetadata     struct {
													WebCommandMetadata struct {
														SendPost bool `json:"sendPost"`
													} `json:"webCommandMetadata"`
												} `json:"commandMetadata"`
												SignalServiceEndpoint struct {
													Signal  string `json:"signal"`
													Actions []struct {
														ClickTrackingParams  string `json:"clickTrackingParams"`
														AddToPlaylistCommand struct {
															OpenMiniplayer      bool   `json:"openMiniplayer"`
															OpenListPanel       bool   `json:"openListPanel"`
															VideoID             string `json:"videoId"`
															ListType            string `json:"listType"`
															OnCreateListCommand struct {
																ClickTrackingParams string `json:"clickTrackingParams"`
																CommandMetadata     struct {
																	WebCommandMetadata struct {
																		SendPost bool   `json:"sendPost"`
																		APIURL   string `json:"apiUrl"`
																	} `json:"webCommandMetadata"`
																} `json:"commandMetadata"`
																CreatePlaylistServiceEndpoint struct {
																	VideoIds []string `json:"videoIds"`
																	Params   string   `json:"params"`
																} `json:"createPlaylistServiceEndpoint"`
															} `json:"onCreateListCommand"`
															VideoIds []string `json:"videoIds"`
														} `json:"addToPlaylistCommand"`
													} `json:"actions"`
												} `json:"signalServiceEndpoint"`
											} `json:"untoggledServiceEndpoint"`
											UntoggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"untoggledAccessibility"`
											ToggledAccessibility struct {
												AccessibilityData struct {
													Label string `json:"label"`
												} `json:"accessibilityData"`
											} `json:"toggledAccessibility"`
											TrackingParams string `json:"trackingParams"`
										} `json:"thumbnailOverlayToggleButtonRenderer,omitempty"`
										ThumbnailOverlayNowPlayingRenderer struct {
											Text struct {
												Runs []struct {
													Text string `json:"text"`
												} `json:"runs"`
											} `json:"text"`
										} `json:"thumbnailOverlayNowPlayingRenderer,omitempty"`
									} `json:"thumbnailOverlays"`
									Accessibility struct {
										AccessibilityData struct {
											Label string `json:"label"`
										} `json:"accessibilityData"`
									} `json:"accessibility"`
								} `json:"compactVideoRenderer,omitempty"`
							} `json:"contents"`
						} `json:"itemSectionRenderer,omitempty"`
					} `json:"results"`
				} `json:"secondaryResults"`
			} `json:"secondaryResults"`
		} `json:"twoColumnWatchNextResults"`
	} `json:"contents"`

	// Chapters
	PlayerOverlays struct {
		PlayerOverlayRenderer struct {
			DecoratedPlayerBarRenderer struct {
				DecoratedPlayerBarRenderer struct {
					PlayerBar struct {
						MultiMarkersPlayerBarRenderer struct {
							MarkersMap []struct {
								Key   string `json:"key"`
								Value struct {
									Chapters []struct {
										ChapterRenderer struct {
											Title struct {
												SimpleText string `json:"simpleText"`
											} `json:"title"`
											TimeRangeStartMillis int `json:"timeRangeStartMillis"`
										} `json:"chapterRenderer"`
									} `json:"chapters"`
								} `json:"value"`
							} `json:"markersMap"`
						} `json:"multiMarkersPlayerBarRenderer"`
					} `json:"playerBar"`
				} `json:"decoratedPlayerBarRenderer"`
			} `json:"decoratedPlayerBarRenderer"`
		} `json:"playerOverlayRenderer"`
	} `json:"playerOverlays"`

	FrameworkUpdates struct {
		EntityBatchUpdate struct {
			Mutations []struct {
				Payload struct {
					SubscriptionStateEntity struct {
						Subscribed bool `json:"subscribed"`
					} `json:"subscriptionStateEntity"`
				} `json:"payload"`
			} `json:"mutations"`
		} `json:"entityBatchUpdate"`
	} `json:"frameworkUpdates"`
}
