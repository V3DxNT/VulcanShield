Right now all the transactions on the / page show as accepted, I want to see it as Challenge -> Accepted/Rejected

And in the top, I can see Valid revenue and then otp accepted, blocked, Fraud Rejected, they are not being accurate and when I click on any, I need to see all the transactions that were under that category, If I click on Challenged, I need to see all the transactions that were challenged, and then the RAG (context) given to the LLM is not being clear, it should see the users past transactions too, it needs to know the last transactions was accepted or rejected!

Isolation Forest is not showing correctly!

I want to see how the outputs are being generated! 

In the extracted evidence signal, it still uses dollar sign instead of INR

Clicking on the XGBoost or Isolation Forest doesn't show up any window to see the reason that they have given these scores 
it says this : High risk because the transaction amount exceeds the customer norm and the policy gate is not satisfied; XGBoost interprets those signals as strong fraud risk.


But I need to see the exact reasons for the transactions have been classified! The formulaes used and how did this number come up! 

I dont think so the ML and the AI services are working properly here! if it does, it needs to know the context, the vector embeddings are not working

Make sure that the services work authentic


___________________________

There are two more problems that the network graph is being cached probably as it is same even on refresh, even on contaiers being restarted
And one feature that you didnt implement was clicking on the "BLOCKED" or "CHALLENGED" should show all the transactions that were blocked or challenged respectively!
One more thing was that the status, is currently
it is Accepted or OTP Rejected
there can be more things that if the ml risk score is above (eg 90) then it is fraud, else if it is low you show accepted, if it is mid, you can show from Challenge -> (outcome of the OTP or RAG Context)
this has to be implemented