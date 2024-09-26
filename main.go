package main

func main() {
	// client, conn := newClient()
	// defer conn.Close()
	//
	// scanner := bufio.NewScanner(os.Stdin)
	// for {
	// 	fmt.Print("input math expression: ")
	// 	if !scanner.Scan() {
	// 		break
	// 	}
	// 	input := scanner.Text()
	// 	if input == "" {
	// 		break
	// 	}
	//
	// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 	defer cancel()
	//
	// 	res, err := client.Calculate(ctx, &pb.Expression{Expression: input})
	// 	if err != nil {
	// 		fmt.Printf("error: %v\n", err)
	// 	} else {
	// 		fmt.Println(res.Result)
	// 	}
	// }
}

// func newClient() (pb.CalculatorClient, *grpc.ClientConn) {
// 	opts := []grpc.DialOption{
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	}
//
// 	conn, err := grpc.NewClient(constant.Address, opts...)
// 	if err != nil {
// 		log.Fatalf("grpc.NewClient failed: %v", err)
// 	}
// 	return pb.NewCalculatorClient(conn), conn
// }
