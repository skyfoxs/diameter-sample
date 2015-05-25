package dictionary

var AppDictionary = `	<?xml version="1.0" encoding="UTF-8"?>
<diameter>

	<application id="4">
		<!-- Huawei Diameter Credit Control Application -->
		<command code="272" short="CC" name="Credit Control">
			<request>
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Destination-Host" required="false" max="1"/>
				<rule avp="Destination-Realm" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="Service-Context-Id" required="true" max="1"/>
				<rule avp="CC-Request-Type" required="true" max="1"/>
				<rule avp="CC-Request-Number" required="true" max="1"/>
				<rule avp="Requested-Action" required="false" max="1"/>
				<rule avp="Event-Timestamp" required="false" max="1"/>
				<rule avp="Service-Identifier" required="false" max="1"/>
				<rule avp="Route-Record" required="false" max="1"/>
				<rule avp="Account-Code" required="false" max="1"/>
				<rule avp="Subscription-Id" required="false" max="1"/>
				<rule avp="Service-Information" required="false" max="1"/>
            </request>
            <answer>
				<rule avp="Session-Id" required="true" max="1"/>
				<rule avp="Result-Code" required="true" max="1"/>
				<rule avp="Origin-Host" required="true" max="1"/>
				<rule avp="Origin-Realm" required="true" max="1"/>
				<rule avp="Auth-Application-Id" required="true" max="1"/>
				<rule avp="CC-Request-Type" required="true" max="1"/>
				<rule avp="CC-Request-Number" required="true" max="1"/>
				<rule avp="Service-Information" required="false" max="1"/>
            </answer>
        </command>

		<avp name="Service-Information" code="873" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<rule avp="Balance-Information" required="false" max="1"/>
				<rule avp="Recharge-Information" required="false" max="1"/>
            </data>
        </avp>

		<avp name="Recharge-Information" code="20800" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<!-- has -->
				<rule avp="Internal-Serial-No" required="false" max="1"/>
				<!-- has -->
				<rule avp="Active-Period" required="false" max="1"/>
				<!-- has -->
				<rule avp="Grace-Period" required="false" max="1"/>
				<!-- has -->
				<rule avp="Disable-Period" required="false" max="1"/>
				<!-- has -->
				<rule avp="New-Balance" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Grade" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Amount" required="false" max="1"/>
				<!-- has -->
				<rule avp="Repay-Amount" required="false" max="1"/>
				<!-- has -->
				<rule avp="ETU-Grace-Period" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Grace-Period" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Acct-Type" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Balance" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Poundage" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-Charge-Info" required="false" max="1"/>
				<!-- has -->
				<rule avp="Service-Status" required="false" max="1"/>
				<!-- has -->
				<rule avp="Original-Loan-Amount" required="false" max="1"/>
				<!-- has -->
				<rule avp="Loan-Time" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Service-Status" code="22308" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned32"/>
		</avp>

		<avp name="Original-Loan-Amount" code="22309" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Loan-Time" code="22310" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Time"/>
		</avp>

		<avp name="Internal-Serial-No" code="22318" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Active-Period" code="20733" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Grace-Period" code="20734" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Disable-Period" code="20735" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="New-Balance" code="22319" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Loan-Grade" code="22300" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Loan-Amount" code="22301" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Repay-Amount" code="22302" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="ETU-Grace-Period" code="22304" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Time"/>
		</avp>

		<avp name="Loan-Grace-Period" code="22305" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Time"/>
		</avp>

		<avp name="Loan-Acct-Type" code="22306" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>
		<avp name="Loan-Balance" code="22307" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Loan-Poundage" code="22315" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Account-Charge-Info" code="20349" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<!-- has -->
				<rule avp="Account-Id" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-Type" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-Type-Desc" required="false" max="1"/>
				<!-- has -->
				<rule avp="Current-Account-Balance" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-Balance-Change" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-Being-Date" required="false" max="1"/>
				<!-- has -->
				<rule avp="Account-End-Date" required="false" max="1"/>
				<!-- has -->
				<rule avp="Measure-Type" required="false" max="1"/>
				<!-- has -->
				<rule avp="Offer-Id" required="false" max="1"/>
				<!-- has not -->
				<rule avp="Time-Schema-Id" required="false" max="1"/>
			</data>
		</avp>

		<avp name="Account-Balance-Change" code="20351" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Account-Being-Date" code="22123" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Account-End-Date" code="20359" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Account-Code" code="30951" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Management-Status" code="22149" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Balance-Information" code="21100" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<rule avp="First-Active-Date" required="false" max="1"/>
				<rule avp="Subscriber-State" required="false" max="1"/>
				<rule avp="Active-Period" required="false" max="1"/>
				<rule avp="Grace-Period" required="false" max="1"/>
				<rule avp="Disable-Period" required="false" max="1"/>
				<rule avp="Period-Balance" required="false" max="1"/>
				<rule avp="Reserve-Amount" required="false" max="1"/>
				<rule avp="Next-Bill-Date" required="false" max="1"/>
				<rule avp="Domestic-Unbilled-Amount1" required="false" max="1"/>
				<rule avp="Domestic-Unbilled-Amount2" required="false" max="1"/>
				<rule avp="IR-Unbilled-Amount1" required="false" max="1"/>
				<rule avp="IR-Unbilled-Amount2" required="false" max="1"/>
				<rule avp="Domestic-Available-Credit" required="false" max="1"/>
				<rule avp="Domestic-Permanent-Credit-Limit" required="false" max="1"/>
				<rule avp="IR-Credit-Limit" required="false" max="1"/>
				<rule avp="Language-IVR" required="false" max="1"/>
				<rule avp="Language-SMS" required="false" max="1"/>
				<rule avp="Language-USSD" required="false" max="1"/>
				<rule avp="Account-Change-Info" required="false" max="1"/>
				<rule avp="Calling-Party-Address" required="false" max="1"/>
				<rule avp="Calling-Cell-Id-Or-SAI" required="false" max="1"/>
				<rule avp="Time-Zone" required="false" max="1"/>
				<rule avp="Access-Method" required="false" max="1"/>
				<rule avp="Account-Query-Method" required="false" max="1"/>
				<rule avp="SSP-Time" required="false" max="1"/>
				<rule avp="Offer-Id-Range" required="false" max="1"/>
			</data>
		</avp>

		<avp name="First-Active-Date" code="20771" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Subscriber-State" code="30814" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned32"/>
		</avp>

		<avp name="Active-Period" code="20733" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Grace-Period" code="20734" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Disable-Period" code="20735" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Prepaid-Balance" code="30841" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Reserve-Amount" code="31800" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Next-Bill-Date" code="31801" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Domestic-Unbilled-Amount1" code="31802" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Domestic-Unbilled-Amount2" code="31803" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="IR-Unbilled-Amount1" code="31804" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="IR-Unbilled-Amount2" code="31805" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Domestic-Available-Credit" code="31806" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Domestic-Permanent-Credit-Limit" code="31807" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="IR-Credit-Limit" code="31808" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Language-IVR" code="21194" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Language-SMS" code="21195" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Language-USSD" code="30939" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Calling-Party-Address" code="20336" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Calling-Cell-Id-Or-SAI" code="20303" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Time-Zone" code="20324" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned32"/>
		</avp>

		<avp name="Access-Method" code="20340" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned32"/>
		</avp>

		<avp name="Account-Query-Method" code="20346" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Unsigned32"/>
		</avp>

		<avp name="SSP-Time" code="20386" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Time"/>
		</avp>

		<avp name="Offer-Id-Range" code="22162" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Account-Change-Info" code="20349" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<rule avp="Account-Id" required="false" max="1"/>
				<rule avp="Account-Type" required="false" max="1"/>
				<rule avp="Account-Type-Desc" required="false" max="1"/>
				<rule avp="Account-Begin-Date" required="false" max="1"/>
				<rule avp="Related-Type" required="false" max="1"/>
				<rule avp="Related-Object-Id" required="false" max="1"/>
				<rule avp="Current-Account-Balance" required="false" max="1"/>
				<rule avp="Account-End-Date" required="false" max="1"/>
				<rule avp="Measure-Typ" required="false" max="1"/>
				<rule avp="Share-Flag" required="false" max="1"/>
            </data>
        </avp>

		<avp name="Account-Id" code="20357" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Account-Type" code="20372" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Account-Type-Desc" code="22320" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Account-Begin-Date" code="22123" must="-" may="P,M" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Related-Type" code="22322" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Related-Object-Id" code="22323" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Current-Account-Balance" code="20350" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer64"/>
		</avp>

		<avp name="Account-End-Date" code="20359" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="OctetString"/>
		</avp>

		<avp name="Measure-Type" code="20353" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Share-Flag" code="30941" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Offer-Information" code="23000" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<rule avp="Offer-Info" required="false" max="1"/>
            </data>
        </avp>

		<avp name="Offer-Info" code="22150" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Grouped">
				<rule avp="Offer-Id" required="false" max="1"/>
				<rule avp="Offer-Order-Key" required="false" max="1"/>
				<rule avp="Effective-Time" required="false" max="1"/>
				<rule avp="Expire-Time" required="false" max="1"/>
				<rule avp="Status" required="false" max="1"/>
				<rule avp="Cur-Cycle-Start-Time" required="false" max="1"/>
				<rule avp="Cur-Cycle-End-Time" required="false" max="1"/>
				<rule avp="Current-Cycle" required="false" max="1"/>
				<rule avp="Total-Cycle" required="false" max="1"/>
				<rule avp="Offer-Order-Integration-Key" required="false" max="1"/>
				<rule avp="External-Offer-Code" required="false" max="1"/>
            </data>
        </avp>

		<avp name="Offer-Id" code="22151" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Offer-Order-Key" code="22152" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Effective-Time" code="22153" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Expire-Time" code="22154" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Status" code="22155" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Cur-Cycle-Start-Time" code="22156" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Cur-Cycle-End-Time" code="22157" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="Current-Cycle" code="22158" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Total-Cycle" code="22159" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="Integer32"/>
		</avp>

		<avp name="Offer-Order-Integration-Key" code="22160" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>

		<avp name="External-Offer-Code" code="22144" must="M" may="P" must-not="V" may-encrypt="Y">
			<data type="UTF8String"/>
		</avp>
    </application>
</diameter>
`
