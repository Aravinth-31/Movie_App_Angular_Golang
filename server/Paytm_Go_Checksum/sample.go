/* More Details: https://developer.paytm.com/docs/checksum/#go */

package main

import (
	"fmt"

	PaytmChecksum "./paytm"
)

func main() {

	/* Generate Checksum via Map */
	/* initialize an map */
	paytmParams := make(map[string]string)

	paytmParams = map[string]string{
		"MID":      "guZEbt55224693629247guZEbt55224693629247",
		"ORDER_ID": "YOUR_ORDER_ID_HERE",
	}

	/**
	* Generate checksum by parameters we have
	* Find your Merchant Key in your Paytm Dashboard at https://dashboard.paytm.com/next/apikeys
	 */
	paytmChecksum := PaytmChecksum.GenerateSignature(paytmParams, "ltARhiHzhjdbU0%K")
	verifyChecksum := PaytmChecksum.VerifySignature(paytmParams, "ltARhiHzhjdbU0%K", paytmChecksum)

	fmt.Printf("GenerateSignature Returns: %s\n", paytmChecksum)
	fmt.Printf("VerifySignature Returns: %t\n\n", verifyChecksum)

	/* Generate Checksum via String */
	/* initialize JSON String */
	body := "{\"mid\":\"guZEbt55224693629247\",\"orderId\":\"YOUR_ORDER_ID_HERE\"}"

	/**
	* Generate checksum by parameters we have
	* Find your Merchant Key in your Paytm Dashboard at https://dashboard.paytm.com/next/apikeys
	 */
	paytmChecksum = PaytmChecksum.GenerateSignatureByString(body, "ltARhiHzhjdbU0%K")
	verifyChecksum = PaytmChecksum.VerifySignatureByString(body, "ltARhiHzhjdbU0%K", paytmChecksum)

	fmt.Printf("GenerateSignatureByString Returns: %s\n", paytmChecksum)
	fmt.Printf("VerifySignatureByString Returns: %t\n\n", verifyChecksum)
}
