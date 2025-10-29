# Common Patterns

Common patterns and best practices for using the SDK.

## Pagination

Handle paginated results:

```go
func getAllProjects(client *pipeops.Client, ctx context.Context) ([]pipeops.Project, error) {
    var allProjects []pipeops.Project
    page := 1
    limit := 100
    
    for {
        projects, _, err := client.Projects.List(ctx, &pipeops.ProjectListOptions{
            Page:  page,
            Limit: limit,
        })
        if err != nil {
            return nil, err
        }
        
        allProjects = append(allProjects, projects.Data.Projects...)
        
        // Check if there are more pages
        if len(projects.Data.Projects) < limit {
            break
        }
        
        page++
    }
    
    return allProjects, nil
}
```

## Bulk Operations

Process multiple resources:

```go
func deployAllProjects(client *pipeops.Client, ctx context.Context, projectUUIDs []string) error {
    for _, uuid := range projectUUIDs {
        _, _, err := client.Projects.Deploy(ctx, uuid)
        if err != nil {
            log.Printf("Failed to deploy %s: %v", uuid, err)
            continue
        }
        
        fmt.Printf("Deployed: %s\n", uuid)
    }
    
    return nil
}
```

## Concurrent Operations

Execute operations concurrently:

```go
func deployProjectsConcurrently(client *pipeops.Client, ctx context.Context, projectUUIDs []string) {
    var wg sync.WaitGroup
    
    for _, uuid := range projectUUIDs {
        wg.Add(1)
        go func(projectUUID string) {
            defer wg.Done()
            
            _, _, err := client.Projects.Deploy(ctx, projectUUID)
            if err != nil {
                log.Printf("Failed to deploy %s: %v", projectUUID, err)
                return
            }
            
            fmt.Printf("Deployed: %s\n", projectUUID)
        }(uuid)
    }
    
    wg.Wait()
}
```

## Resource Polling

Poll for resource state changes:

```go
func waitForDeployment(client *pipeops.Client, ctx context.Context, projectUUID string) error {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    timeout := time.After(10 * time.Minute)
    
    for {
        select {
        case <-timeout:
            return fmt.Errorf("deployment timeout")
            
        case <-ticker.C:
            project, _, err := client.Projects.Get(ctx, projectUUID)
            if err != nil {
                return err
            }
            
            status := project.Data.Project.Status
            if status == "deployed" {
                return nil
            } else if status == "failed" {
                return fmt.Errorf("deployment failed")
            }
            
            fmt.Printf("Status: %s\n", status)
        }
    }
}
```

## Caching

Cache frequently accessed data:

```go
type CachedClient struct {
    client *pipeops.Client
    cache  map[string]interface{}
    mu     sync.RWMutex
}

func (c *CachedClient) GetProject(ctx context.Context, uuid string) (*pipeops.Project, error) {
    c.mu.RLock()
    if cached, ok := c.cache[uuid]; ok {
        c.mu.RUnlock()
        return cached.(*pipeops.Project), nil
    }
    c.mu.RUnlock()
    
    project, _, err := c.client.Projects.Get(ctx, uuid)
    if err != nil {
        return nil, err
    }
    
    c.mu.Lock()
    c.cache[uuid] = &project.Data.Project
    c.mu.Unlock()
    
    return &project.Data.Project, nil
}
```
