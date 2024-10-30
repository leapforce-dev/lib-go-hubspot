package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type BlogPost struct {
	Id   string `json:"id"`
	Slug string `json:"slug"`
	Url  string `json:"url"`
}

type BlogPostsResponse struct {
	Results []BlogPost `json:"results"`
	Paging  *Paging    `json:"paging"`
}

type GetBlogsConfig struct {
	Limit *uint
	After *string
}

// GetBlogPosts returns all blog posts
func (service *Service) GetBlogPosts(config *GetBlogsConfig) (*[]BlogPost, *errortools.Error) {
	values := url.Values{}
	endpoint := "blogs/posts"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
	}

	after := ""

	if config != nil {
		if config.After != nil {
			after = *config.After
		}
	}

	var blogPosts []BlogPost

	for {
		blogPostsResponse := BlogPostsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCms(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &blogPostsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		blogPosts = append(blogPosts, blogPostsResponse.Results...)

		if config != nil {
			if config.After != nil { // explicit after parameter requested
				break
			}
		}

		if blogPostsResponse.Paging == nil {
			break
		}

		if blogPostsResponse.Paging.Next.After == "" {
			break
		}

		after = blogPostsResponse.Paging.Next.After
	}

	return &blogPosts, nil
}

// DeleteBlogPost deletes a specific blog post
func (service *Service) DeleteBlogPost(id string) *errortools.Error {
	endpoint := "blogs/posts"

	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlCms(fmt.Sprintf("%s/%s", endpoint, id)),
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}
